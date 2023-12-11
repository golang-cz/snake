package main

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-cz/snake/proto"
)

func (s *Server) Run(ctx context.Context) error {
	go s.generateFood()
	go s.generatePlayers(3)
	go s.generateSnakeTurns()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.gameTick(); err != nil {
				return fmt.Errorf("advancing the game: %w", err)
			}
		}
	}

	return nil
}

func (s *Server) gameTick() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// snakeGrid := make([][]*proto.Snake, s.state.Width)
	// for x := range snakeGrid {
	// 	snakeGrid[x] = make([]*proto.Snake, s.state.Height)
	// }
	// for _, snake := range s.state.Snakes {
	// 	for _, square := range snake.Body {
	// 		snakeGrid[square.X][square.Y] = snake
	// 	}
	// }

	//	snakeGrid := make([][]*proto.Snake, s.state.Width)
	// for x := range snakeGrid {
	// 	snakeGrid[x] = make([]*proto.Snake, s.state.Height)
	// }
	// for _, snake := range s.state.Snakes {
	// 	for _, square := range snake.Body {
	// 		snakeGrid[square.X][square.Y] = snake
	// 	}
	// }
	// for _, item := range s.state.Items {
	// 	itemGrid[item.Body.X][item.Body.Y] = item
	// }

	// Move snakes.
	for _, snake := range s.state.Snakes {
		newSquare := &proto.Square{}

		if len(snake.NextDirections) > 0 {
			snake.Direction = snake.NextDirections[0]
			snake.NextDirections = snake.NextDirections[1:]
		}

		switch *snake.Direction {
		case proto.Direction_up:
			newSquare.X = snake.Body[0].X
			newSquare.Y = min(snake.Body[0].Y-1, s.state.Height-1)
		case proto.Direction_down:
			newSquare.X = snake.Body[0].X
			newSquare.Y = min(snake.Body[0].Y+1) % s.state.Height
		case proto.Direction_left:
			newSquare.X = min(snake.Body[0].X-1, s.state.Width-1)
			newSquare.Y = snake.Body[0].Y
		case proto.Direction_right:
			newSquare.X = min(snake.Body[0].X+1) % s.state.Width
			newSquare.Y = snake.Body[0].Y
		}

		// Look through items.. TODO: map[Square]*Item?
		eat := false
		for i, item := range s.state.Items {
			if item.Body.X == newSquare.X && item.Body.Y == newSquare.Y {
				eat = true
				delete(s.state.Items, i)
				break
			}
		}

		if eat {
			snake.Body = append([]*proto.Square{newSquare}, snake.Body...)
		} else {
			snake.Body = append([]*proto.Square{newSquare}, snake.Body[:len(snake.Body)-1]...)
		}
	}

	return s.sendState(s.state)
}

func (s *Server) eventLoop(ctx context.Context) error {
	for event := range s.events {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := s.updateState(event); err != nil {
			return fmt.Errorf("updating state: %w", err)
		}
	}

	return nil
}

func (s *Server) updateState(events ...*proto.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// for _, event := range events {
	// 	// switch event.Type {

	// 	// }
	// }
	return nil
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
