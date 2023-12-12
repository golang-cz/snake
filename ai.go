package main

import (
	"image"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/fzipp/astar"
	"github.com/golang-cz/snake/proto"
)

func (s *Server) createFood() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state.Items[s.lastItemId] = &proto.Item{
		Id:    s.lastItemId,
		Color: "red",
		Body: &proto.Square{
			X: uint(rand.Intn(int(s.state.Width))),
			Y: uint(rand.Intn(int(s.state.Height))),
		},
	}
	s.lastItemId++
}

func (s *Server) generateFood() {
	s.createFood()
	s.createFood()
	s.createFood()
	for {
		<-time.After(2 * time.Second)
		s.createFood()
	}
}

func (s *Server) currentGrid() grid {
	grid := make(grid, s.state.Height)
	for y := 0; y < int(s.state.Height); y++ {
		grid[y] = strings.Repeat(" ", int(s.state.Width))
	}

	for _, item := range s.state.Items {
		grid.put(image.Point{X: int(item.Body.X), Y: int(item.Body.Y)}, '*')
	}

	for _, snake := range s.state.Snakes {
		for _, square := range snake.Body {
			grid.put(image.Point{X: int(square.X), Y: int(square.Y)}, snakeRune(snake.Id))
		}
	}

	return grid
}

func (s *Server) generateSnakeTurns(grid grid) {
	// Turn "AI" snakes to the closes food using A* algorithm.
	for _, snake := range s.state.Snakes {
		if snake.Name != "AI" {
			continue
		}

		// TODO: Make this based on speed and levels?
		// badLuck := rand.Intn(10)
		// if badLuck>7 {
		// 	continue
		// }

		snakeHead := squareToPoint(snake.Body[0])

		// Find closest food.
		var closestFood astar.Path[image.Point]
		shortestPathLen := math.MaxInt

		for _, item := range s.state.Items {
			food := squareToPoint(item.Body)

			path := astar.FindPath[image.Point](grid, snakeHead, food, distance, distance)
			if len(path) > 1 && len(path) < shortestPathLen {
				closestFood = path[1:]
				shortestPathLen = len(closestFood)
			}
		}

		if len(closestFood) > 0 {
			// Mark shortest path with dots.
			// for _, p := range closestFood[:len(closestFood)-1] {
			// 	if grid.get(p) == ' ' {
			// 		grid.put(p, '.')
			// 	}
			// }

			// Turn snake to the direction of the closest food.
			nextSquare := closestFood[0]

			switch {
			case snakeHead.X < nextSquare.X:
				if turnSnake(snake, &right, 0) == proto.ErrTurnAbout {
					turnSnake(snake, &up, 0)
				}
			case snakeHead.X > nextSquare.X:
				if turnSnake(snake, &left, 0) == proto.ErrTurnAbout {
					turnSnake(snake, &down, 0)
				}
			case snakeHead.Y < nextSquare.Y:
				if turnSnake(snake, &down, 0) == proto.ErrTurnAbout {
					turnSnake(snake, &left, 0)
				}
			case snakeHead.Y > nextSquare.Y:
				if turnSnake(snake, &up, 0) == proto.ErrTurnAbout {
					turnSnake(snake, &right, 0)
				}
			}
		}
	}
}
