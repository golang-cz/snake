package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/golang-cz/snake/proto"
)

func (s *Server) generateFood() {
	// Generate new food.
	for {
		<-time.After(5 * time.Second)
		func() {
			s.mu.Lock()
			defer s.mu.Unlock()

			s.state.Items[s.lastItemId] = &proto.Item{
				Id:    s.lastItemId,
				Color: "red",
				Body: &proto.Square{
					X: uint(rand.Intn(int(s.state.Width))),
					Y: uint(rand.Intn(int(s.state.Height))),
				},
			}
			s.lastItemId++
		}()
	}
}

func (s *Server) generatePlayers(num int) {
	return
	// Simulate new players joining the game.
	for i := 0; i < num; i++ {
		<-time.After(3 * time.Second)
		s.CreateSnake(context.Background(), "AI")
	}
}

func (s *Server) generateSnakeTurns() {
	// Simulate players turning.
	for {
		if snake, ok := s.state.Snakes[uint64(rand.Int63n(int64(s.lastSnakeId)))]; ok {
			if snake.Name != "AI" {
				continue
			}
			s.TurnSnake(context.Background(), snake.Id, randDirection())
		}
		<-time.After(100 * time.Millisecond)
	}
}
