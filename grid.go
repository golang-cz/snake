package main

import (
	"fmt"
	"image"
	"math"
	"strings"

	"github.com/golang-cz/snake/proto"
)

type grid []string

func (g grid) String() string {
	var b strings.Builder
	b.WriteByte('\n')
	b.WriteString(strings.Repeat("▮", len(g)+2))
	b.WriteByte('\n')
	for _, row := range g {
		b.WriteString("▮")
		b.WriteString(row)
		b.WriteString("▮")
		b.WriteByte('\n')
	}
	b.WriteString(strings.Repeat("▮", len(g)+2))
	b.WriteByte('\n')
	return b.String()
}

// Neighbours implements the astar.Graph interface
func (g grid) Neighbours(p image.Point) []image.Point {
	offsets := []image.Point{
		image.Pt(0, -1), // North
		image.Pt(1, 0),  // East
		image.Pt(0, 1),  // South
		image.Pt(-1, 0), // West
	}
	res := make([]image.Point, 0, 4)
	for _, off := range offsets {
		q := p.Add(off)
		if g.isFreeAt(q) {
			res = append(res, q)
		}
	}
	return res
}

func (g grid) isFreeAt(p image.Point) bool {
	return g.isInBounds(p) && g[p.Y][p.X] != '#'
}

func (g grid) isInBounds(p image.Point) bool {
	return p.Y >= 0 && p.X >= 0 && p.Y < len(g) && p.X < len(g[p.Y])
}

func (g grid) put(p image.Point, c rune) {
	g[p.Y] = g[p.Y][:p.X] + string(c) + g[p.Y][p.X+1:]
}

func (g grid) print() {
	for _, row := range g {
		fmt.Println(row)
	}
}

func distance(p, q image.Point) float64 {
	d := q.Sub(p)
	return math.Sqrt(float64(d.X*d.X + d.Y*d.Y))
}

func SquareToPoint(s *proto.Square) image.Point {
	return image.Point{
		X: int(s.X),
		Y: int(s.Y),
	}
}
