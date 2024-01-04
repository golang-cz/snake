package main

import (
	"context"
	"fmt"

	"github.com/golang-cz/snake/proto"
)

func (s *Server) JoinGame(ctx context.Context, stream proto.JoinGameStreamWriter) error {
	events := make(chan *proto.Update, 10)

	update, subscriptionId := s.subscribe(events)
	defer s.unsubscribe(subscriptionId)

	// Send initial state.
	if err := stream.Write(update, nil); err != nil {
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

		case update := <-events:
			if err := stream.Write(update, nil); err != nil {
				return err
			}
		}
	}
}

func (s *Server) CreateSnake(ctx context.Context, username string) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.createSnake(username)
}

func (s *Server) TurnSnake(ctx context.Context, snakeId uint64, direction *proto.Direction) error {
	if direction == nil {
		return fmt.Errorf("nil direction")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	snake, ok := s.state.Snakes[snakeId]
	if !ok {
		return proto.ErrSnakeNotFound.WithCause(fmt.Errorf("snakeId %v not found", snakeId))
	}

	// Turn snake, if possible, and buffer up to 2 actions.
	return turnSnake(snake, direction, 2)
}
