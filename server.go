package main

import (
	"sync"

	"github.com/golang-cz/snake/proto"
)

type Server struct {
	mu     sync.Mutex
	state  *proto.State
	events chan *proto.Event
	subs   map[uint64]chan *proto.State

	lastSnakeId uint64
	lastItemId  uint64
	lastSubId   uint64
}

var (
	down  = proto.Direction_down
	up    = proto.Direction_up
	left  = proto.Direction_left
	right = proto.Direction_right
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
		subs:   map[uint64]chan *proto.State{},
	}
}
