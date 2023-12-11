package main

import (
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
