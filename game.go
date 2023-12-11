package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-cz/snake/proto"
)

func (s *Server) Run(ctx context.Context) error {
	go s.generateFood()
	go s.generatePlayers(3)
	go s.generateSnakeTurns()

	ticker := time.NewTicker(100 * time.Millisecond)
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

	// Move snakes.
	for _, snake := range s.state.Snakes {
		newSquare := &proto.Square{}

		switch snake.Direction {
		case &up:
			newSquare.X = snake.Body[0].X
			newSquare.Y = min(snake.Body[0].Y-1, s.state.Height)
		case &down:
			newSquare.X = snake.Body[0].X
			newSquare.Y = min(snake.Body[0].Y+1) % s.state.Height
		case &left:
			newSquare.X = min(snake.Body[0].X-1, s.state.Width)
			newSquare.Y = snake.Body[0].Y
		case &right:
			newSquare.X = min(snake.Body[0].X+1) % s.state.Width
			newSquare.Y = snake.Body[0].Y
		}

		log.Print("%#v", newSquare)

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
