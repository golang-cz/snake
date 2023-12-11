package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang-cz/snake/proto"
)

type Server struct {
	mu        sync.Mutex
	state     *proto.State
	events    chan *proto.Event
	subs      map[uint64]chan *proto.State
	lastSubId uint64
}

func NewSnakeServer() *Server {
	return &Server{
		state:  &proto.State{},
		events: make(chan *proto.Event, 100000),
		subs:   map[uint64]chan *proto.State{},
	}
}

func (s *Server) runGame(ctx context.Context) error {
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
