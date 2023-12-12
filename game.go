package main

import (
	"context"
	"fmt"
	"image"
	"time"

	"github.com/golang-cz/snake/proto"
)

func (s *Server) Run(ctx context.Context) error {
	go s.generateFood()

	for i := 0; i < 3; i++ {
		s.createSnake("AI")
	}

	ticker := time.NewTicker(75 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.gameTick(); err != nil {
				return fmt.Errorf("advancing the game: %w", err)
			}
		}
	}
}

func (s *Server) gameTick() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	grid := s.currentGrid()
	defer fmt.Println(grid)

	s.generateSnakeTurns(grid)

	// Move snakes.
	for _, snake := range s.state.Snakes {
		next := &proto.Square{}

		if len(snake.NextDirections) > 0 {
			snake.Direction = snake.NextDirections[0]
			snake.NextDirections = snake.NextDirections[1:]
		}

		switch *snake.Direction {
		case proto.Direction_up:
			next.X = snake.Body[0].X
			next.Y = min(snake.Body[0].Y-1, s.state.Height-1)
		case proto.Direction_down:
			next.X = snake.Body[0].X
			next.Y = min(snake.Body[0].Y+1) % s.state.Height
		case proto.Direction_left:
			next.X = min(snake.Body[0].X-1, s.state.Width-1)
			next.Y = snake.Body[0].Y
		case proto.Direction_right:
			next.X = min(snake.Body[0].X+1) % s.state.Width
			next.Y = snake.Body[0].Y
		}

		nextPoint := image.Point{X: int(next.X), Y: int(next.Y)}

		switch grid.get(squareToPoint(next)) {
		case '*':
			// Move snake's head.
			grid.put(nextPoint, snakeRune(snake.Id))

			snake.Body = append([]*proto.Square{next}, snake.Body...)

			// Remove food.
			for _, item := range s.state.Items {
				if item.Body.X == next.X && item.Body.Y == next.Y {
					delete(s.state.Items, item.Id)
				}
			}

		case ' ', '.', snakeRune(snake.Id):
			// Move snake's head & remove tail.
			tailPoint := squareToPoint(snake.Body[len(snake.Body)-1])
			grid.put(tailPoint, ' ')
			grid.put(nextPoint, snakeRune(snake.Id))

			snake.Body = append([]*proto.Square{next}, snake.Body[:len(snake.Body)-1]...)

		default:
			// Crashed into another snake.
			delete(s.state.Snakes, snake.Id)

			// Create food from snake's body.
			for i, square := range snake.Body {
				if i%2 == 0 {
					s.state.Items[s.lastItemId] = &proto.Item{
						Id:    s.lastItemId,
						Color: "red",
						Body:  square,
					}
					s.lastItemId++
				}
			}

			// Reborn AI.
			if snake.Name == "AI" {
				go func() {
					<-time.After(10 * time.Second)
					s.createSnake("AI")
				}()
			}

		}
	}

	return s.sendState(s.state)
}

// TODO: We send the whole state on each update. Optimize to send events (diffs) only.
func (s *Server) sendState(state *proto.State) error {
	for _, sub := range s.subs {
		sub := sub
		go func() {
			sub <- state
		}()
	}
	return nil
}

var runes = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

func snakeRune(snakeId uint64) rune {
	return runes[int(snakeId)%len(runes)]
}
