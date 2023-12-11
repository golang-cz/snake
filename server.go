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

var (
	down  = proto.Direction_down
	up    = proto.Direction_up
	left  = proto.Direction_left
	right = proto.Direction_right
)

func NewSnakeServer() *Server {
	return &Server{
		state: &proto.State{
			Snakes: []*proto.Snake{
				{
					Id:    1,
					Color: "blue",
					Body: []*proto.Coord{
						{X: 10, Y: 20},
						{X: 10, Y: 21},
						{X: 11, Y: 21},
						{X: 12, Y: 21},
						{X: 13, Y: 21},
						{X: 14, Y: 21},
					},
					Direction: &up,
				},
				{
					Id:    2,
					Color: "green",
					Body: []*proto.Coord{
						{X: 50, Y: 41},
						{X: 50, Y: 42},
						{X: 50, Y: 43},
						{X: 50, Y: 44},
					},
					Direction: &down,
				},
			},
			Items: []*proto.Item{
				{
					Id:    1,
					Color: "red",
					Body: []*proto.Coord{
						{X: 10, Y: 20},
						{X: 10, Y: 21},
						{X: 11, Y: 21},
						{X: 12, Y: 21},
						{X: 13, Y: 21},
						{X: 14, Y: 21},
					},
				},
			},
		},
		events: make(chan *proto.Event, 100000),
		subs:   map[uint64]chan *proto.State{},
	}
}
