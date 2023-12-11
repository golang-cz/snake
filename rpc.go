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

	s.lastSnakeId++

	snakeId := s.lastSnakeId
	s.state.Snakes[snakeId] = &proto.Snake{
		Id:    snakeId,
		Name:  username,
		Color: randColor(),
		Body: []*proto.Square{
			{X: 36, Y: 35},
			{X: 35, Y: 35},
			{X: 34, Y: 35},
		},
		Direction: &right,
	}

	return snakeId, nil
}

func (s *Server) TurnSnake(ctx context.Context, snakeId uint64, direction *proto.Direction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if direction == nil {
		return fmt.Errorf("nil direction")
	}

	snake, ok := s.state.Snakes[snakeId]
	if !ok {
		return proto.ErrSnakeNotFound.WithCause(fmt.Errorf("snakeId %v not found", snakeId))
	}

	lastDirection := *snake.Direction
	if len(snake.NextDirections) > 0 {
		lastDirection = *snake.NextDirections[len(snake.NextDirections)-1]
	}

	// Same direction.
	if lastDirection == *direction {
		return proto.ErrInvalidTurn
	}

	// Disallow turnabouts.
	switch {
	case lastDirection == proto.Direction_up && *direction == proto.Direction_down:
		return proto.ErrInvalidTurn
	case lastDirection == proto.Direction_down && *direction == proto.Direction_up:
		return proto.ErrInvalidTurn
	case lastDirection == proto.Direction_left && *direction == proto.Direction_right:
		return proto.ErrInvalidTurn
	case lastDirection == proto.Direction_right && *direction == proto.Direction_left:
		return proto.ErrInvalidTurn
	}

	if len(snake.NextDirections) > 2 {
		snake.NextDirections = append(snake.NextDirections[:2], direction)
	} else {
		snake.NextDirections = append(snake.NextDirections, direction)
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
