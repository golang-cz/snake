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
			Width:  70,
			Height: 70,
			Snakes: []*proto.Snake{
				{
					Id:    1,
					Color: "blue",
					Body: []*proto.Square{
						{X: 10, Y: 20},
						{X: 10, Y: 21},
						{X: 11, Y: 21},
						{X: 12, Y: 21},
						{X: 13, Y: 21},
						{X: 14, Y: 21},
						{X: 14, Y: 22},
						{X: 14, Y: 23},
						{X: 14, Y: 24},
						{X: 14, Y: 25},
						{X: 13, Y: 25},
						{X: 12, Y: 25},
						{X: 11, Y: 25},
						{X: 11, Y: 26},
						{X: 11, Y: 27},
						{X: 11, Y: 28},
						{X: 11, Y: 29},
					},
					Direction: &up,
				},
				{
					Id:    2,
					Color: "green",
					Body: []*proto.Square{
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
					Body:  &proto.Square{X: 10, Y: 2},
				},
				{
					Id:    2,
					Color: "red",
					Body:  &proto.Square{X: 24, Y: 32},
				},
				{
					Id:    3,
					Color: "red",
					Body:  &proto.Square{X: 50, Y: 49},
				},
				{
					Id:    3,
					Color: "red",
					Body:  &proto.Square{X: 50, Y: 51},
				},
				{
					Id:    3,
					Color: "red",
					Body:  &proto.Square{X: 50, Y: 53},
				},
			},
		},
		events: make(chan *proto.Event, 100000),
		subs:   map[uint64]chan *proto.State{},
	}
}
