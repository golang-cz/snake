package main

import (
	"context"
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/golang-cz/snake/proto"
)

func (s *Server) Run(ctx context.Context) error {
	go s.generateFood()
	go s.generateAISnakes()

	ticker := time.NewTicker(GameTickTime)
	defer ticker.Stop()

	for range ticker.C {
		// if err := s.gameTick(); err != nil {
		// 	return fmt.Errorf("advancing the game: %w", err)
		// }

		s.gameTick()
	}

	return nil
}

func (s *Server) generateAISnakes() {
	for i := 0; i < NumOfAISnakes; i++ {
		time.Sleep(time.Second)
		go s.createSnake(fmt.Sprintf("AI%v", i))
	}
}

func (s *Server) gameTick() {
	s.mu.Lock()
	defer s.mu.Unlock()

	grid := s.currentGrid()
	defer fmt.Println(grid)

	s.generateSnakeTurns(grid)

	u := &proto.Update{
		Diffs: []*proto.Diff{},
		// State: s.state,
	}

	// Move snakes.
	for _, snake := range s.state.Snakes {
		next := &proto.Coordinate{}

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

			snake.Body = append([]*proto.Coordinate{next}, snake.Body...)

			// Remove food.
			for _, item := range s.state.Items {
				if item.Coordinate.X == next.X && item.Coordinate.Y == next.Y {
					delete(s.state.Items, item.Id)
				}
			}

			diff := &proto.Diff{
				X:     next.X,
				Y:     next.Y,
				Color: snake.Color,
				Add:   true,
			}

			u.Diffs = append(u.Diffs, diff)

			snake.Length = len(snake.Body)

		case ' ', '.', snakeRune(snake.Id):
			// Move snake's head & remove tail.
			tail := snake.Body[len(snake.Body)-1]
			tailPoint := squareToPoint(tail)
			grid.put(tailPoint, ' ')
			grid.put(nextPoint, snakeRune(snake.Id))

			headDiff := &proto.Diff{
				X:     next.X,
				Y:     next.Y,
				Color: snake.Color,
				Add:   true,
			}

			tailDiff := &proto.Diff{
				X:   tail.X,
				Y:   tail.Y,
				Add: false,
			}

			u.Diffs = append(u.Diffs, headDiff)
			u.Diffs = append(u.Diffs, tailDiff)

			snake.Body = append([]*proto.Coordinate{next}, snake.Body[:len(snake.Body)-1]...)

		default:
			// Crashed into another snake.
			name := snake.Name
			delete(s.state.Snakes, snake.Id)

			bite := proto.ItemType_bite
			for i, bodyPart := range snake.Body {
				diff := &proto.Diff{
					X: bodyPart.X,
					Y: bodyPart.Y,
				}

				// Create food from snake's body.
				if i%FoodFromDeadSnake == 0 {
					item := &proto.Item{
						Id:         s.lastItemId,
						Color:      "red",
						Coordinate: bodyPart,
						Type:       &bite,
					}

					s.state.Items[s.lastItemId] = item
					s.lastItemId++

					diff.Color = item.Color
					diff.Add = true
				}

				u.Diffs = append(u.Diffs, diff)
			}

			// Reborn AI.
			if strings.Contains(snake.Name, "AI") {
				go func() {
					<-time.After(AISnakeRespawnTime)
					s.createSnake(name)
				}()
			}
		}
	}

	// Generate bite, if snakes ate all bites
	// bites := 0
	// for i := uint64(0); int(i) < len(s.state.Items) && bites == 0; i++ {
	// 	if s.state.Items[i].Type != nil && *s.state.Items[i].Type == proto.ItemType_bite {
	// 		bites++
	// 	}
	// }
	//
	// if bites == 0 {
	// 	go s.createFood()
	// }
	//
	if len(u.Diffs) != 0 {
		s.sendUpdate(u)
	}
}

func (s *Server) sendUpdate(state *proto.Update) {
	for _, sub := range s.subs {
		sub := sub
		go func() {
			sub <- state
		}()
	}
}

// Get letter from A to Z
func snakeRune(snakeId uint64) rune {
	return rune((snakeId % 26) + 65)
}
