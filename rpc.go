package main

import (
	"context"

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
	return 0, nil
}

func (s *Server) TurnSnake(ctx context.Context, snakeId uint64, direction *proto.Direction) error {
	return nil
}
