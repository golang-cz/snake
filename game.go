package main

import (
	"context"
	"fmt"

	"github.com/golang-cz/snake/proto"
)

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
