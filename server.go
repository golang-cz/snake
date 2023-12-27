package main

import (
	"sync"

	"github.com/golang-cz/snake/proto"
)

type Server struct {
	mu     sync.Mutex
	state  *proto.State
	events chan *proto.Event
	subs   map[uint64]chan *proto.Update

	lastSnakeId uint64
	lastItemId  uint64
	lastSubId   uint64
}

var (
	Left  = proto.Direction_left
	Right = proto.Direction_right
	Up    = proto.Direction_up
	Down  = proto.Direction_down
)

func NewSnakeServer() *Server {
	return &Server{
		state: &proto.State{
			Width:  70,
			Height: 70,
			Snakes: map[uint64]*proto.Snake{},
			Items:  map[uint64]*proto.Item{},
		},
		events: make(chan *proto.Event, 100000),
		subs:   map[uint64]chan *proto.Update{},
	}
}
