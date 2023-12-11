package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang-cz/snake/proto"
)

type GameServer struct {
	mu        sync.Mutex
	state     *proto.State
	events    chan *proto.Event
	subs      map[uint64]chan *proto.State
	lastSubId uint64
}

func NewChatServer() *GameServer {
	return &GameServer{
		state:  &proto.State{},
		events: make(chan *proto.Event, 100000),
		subs:   map[uint64]chan *proto.State{},
	}
}

func (s *GameServer) runGame(ctx context.Context) error {
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

func (s *GameServer) updateState(events ...*proto.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// for _, event := range events {
	// 	// switch event.Type {

	// 	// }
	// }
	return nil
}

// TODO: We send the whole state on each update. Optimize to send events (diffs) only.
func (s *GameServer) sendState(state *proto.State) error {
	for _, sub := range s.subs {
		sub := sub
		go func() {
			sub <- state
		}()
	}
	return nil
}

func (s *GameServer) JoinGame(ctx context.Context, stream proto.JoinGameStreamWriter) error {
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

func (s *GameServer) CreateSnake(ctx context.Context, username string) (uint64, error) {
	return 0, nil
}

func (s *GameServer) TurnSnake(ctx context.Context, snakeId uint64, direction *proto.Direction) error {
	return nil
}

func (s *GameServer) subscribe(c chan *proto.State) (*proto.State, uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.lastSubId
	s.subs[id] = c
	s.lastSubId++

	return s.state, id
}

func (s *GameServer) unsubscribe(subscriptionId uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.subs, subscriptionId)
}
