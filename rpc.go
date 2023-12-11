package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/golang-cz/snake/proto"
)

func (s *Server) JoinGame(ctx context.Context, stream proto.JoinGameStreamWriter) error {
	events := make(chan *proto.State, 10)

	state, subscriptionId := s.subscribe(events)
	defer s.unsubscribe(subscriptionId)

	// Send initial state.
	if err := stream.Write(state, nil); err != nil {
		return err
	}

	// Send updates. TODO: Send diffs only.
	for {
		select {
		case <-ctx.Done():
			switch err := ctx.Err(); err {
			case context.Canceled:
				return proto.ErrWebrpcClientDisconnected
			default:
				return proto.ErrWebrpcInternalError
			}

		case state := <-events:
			if err := stream.Write(state, nil); err != nil {
				return err
			}
		}
	}
}

func (s *Server) CreateSnake(ctx context.Context, username string) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state.Snakes[s.lastSnakeId] = &proto.Snake{
		Id:    s.lastSnakeId,
		Color: randColor(),
		Body: []*proto.Square{
			{X: 34, Y: 35},
			{X: 35, Y: 35},
			{X: 36, Y: 35},
		},
		Direction: &right,
	}
	s.lastSnakeId++

	return s.lastSnakeId, nil
}

func (s *Server) TurnSnake(ctx context.Context, snakeId uint64, direction *proto.Direction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if direction == nil {
		return fmt.Errorf("nil direction")
	}

	if snake, ok := s.state.Snakes[snakeId]; ok {
		// Disallow back turns.
		switch {
		case snake.Direction == &up && direction == &down:
			return nil
		case snake.Direction == &down && direction == &up:
			return nil
		case snake.Direction == &left && direction == &right:
			return nil
		case snake.Direction == &right && direction == &left:
			return nil
		}

		snake.Direction = direction
	}

	return nil
}

func randColor() string {
	colors := []string{"blue", "green", "lightgreen", "darkgreen", "lightblue", "darkblue", "pink", "brown", "yellow", "orange", "gray", "lightgray", "purple", "magenta", "black", "aqua"}
	return colors[rand.Intn(len(colors))]
}

func randDirection() *proto.Direction {
	directions := []*proto.Direction{&left, &up, &right, &down}
	return directions[rand.Intn(len(directions))]
}
