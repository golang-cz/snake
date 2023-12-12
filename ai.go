package main

import (
	"image"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/fzipp/astar"
	"github.com/golang-cz/snake/proto"
)

func (s *Server) generateFood() {
	for {
		<-time.After(2 * time.Second)
		func() {
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
		}()
	}
}

func (s *Server) generateSnakeTurns() {
	grid := make(grid, s.state.Height)
	for y := 0; y < int(s.state.Height); y++ {
		grid[y] = strings.Repeat(" ", int(s.state.Width))
	}

	defer log.Println(grid)

	for _, item := range s.state.Items {
		grid.put(image.Point{X: int(item.Body.X), Y: int(item.Body.Y)}, '*')
	}

	for _, snake := range s.state.Snakes {
		for i, square := range snake.Body {
			if i == 0 {
				grid.put(image.Point{X: int(square.X), Y: int(square.Y)}, 'H')
			} else {
				grid.put(image.Point{X: int(square.X), Y: int(square.Y)}, '#')
			}
		}
	}

	// Turn "AI" snakes to the closes food using A* algorithm.
	for _, snake := range s.state.Snakes {
		if snake.Name != "AI" {
			continue
		}

		snakeHead := SquareToPoint(snake.Body[0])

		// Find closest food.
		var closestFood astar.Path[image.Point]
		shortestPathLen := math.MaxInt

		for _, item := range s.state.Items {
			food := SquareToPoint(item.Body)

			path := astar.FindPath[image.Point](grid, snakeHead, food, distance, distance)
			if len(path) > 1 && len(path) < shortestPathLen {
				closestFood = path[1:]
				shortestPathLen = len(closestFood)
			}
		}

		if len(closestFood) > 0 {
			// Mark shortest path with dots.
			for _, p := range closestFood[:len(closestFood)-1] {
				grid.put(p, '.')
			}

			// Turn snake to the direction of the closest food.
			nextSquare := closestFood[0]

			switch {
			case snakeHead.X < nextSquare.X:
				turnSnake(snake, &right, 0)
			case snakeHead.X > nextSquare.X:
				turnSnake(snake, &left, 0)
			case snakeHead.Y < nextSquare.Y:
				turnSnake(snake, &down, 0)
			case snakeHead.Y > nextSquare.Y:
				turnSnake(snake, &up, 0)
			}
		}
	}
}
