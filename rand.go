package main

import (
	"math/rand"

	"github.com/golang-cz/snake/proto"
)

func randColor() string {
	colors := []string{"blue", "green", "lightgreen", "darkgreen", "lightblue", "darkblue", "pink", "brown", "yellow", "orange", "gray", "lightgray", "purple", "magenta", "black", "aqua"}
	return colors[rand.Intn(len(colors))]
}

func randDirection() *proto.Direction {
	directions := []*proto.Direction{&Left, &Up, &Right, &Down}
	return directions[rand.Intn(len(directions))]
}
