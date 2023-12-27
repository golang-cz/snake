package main

import (
	"math/rand"

	"github.com/golang-cz/snake/proto"
)

func (s *Server) createSnake(username string) (uint64, error) {
	randOffset := uint(rand.Intn(10) - rand.Intn(10))

	snakeId := s.lastSnakeId
	snake := &proto.Snake{
		Id:    snakeId,
		Name:  username,
		Color: randColor(),
		Body: []*proto.Coordinate{
			{X: 36, Y: 35 + randOffset},
			{X: 35, Y: 35 + randOffset},
			{X: 34, Y: 35 + randOffset},
		},
		Direction: &Right,
	}

	s.state.Snakes[snakeId] = snake

	u := &proto.Update{
		Diffs: generateDiffFromBody(snake),
	}

	s.sendUpdate(u)

	s.lastSnakeId++

	return snakeId, nil
}

func turnSnake(snake *proto.Snake, direction *proto.Direction, buf int) error {
	lastDirection := *snake.Direction
	if len(snake.NextDirections) > 0 {
		lastDirection = *snake.NextDirections[len(snake.NextDirections)-1]
	}

	// Same direction.
	if lastDirection == *direction {
		return proto.ErrTurnSameDirection
	}

	// Disallow turnabouts.
	switch {
	case lastDirection == proto.Direction_up && *direction == proto.Direction_down:
		return proto.ErrTurnAbout
	case lastDirection == proto.Direction_down && *direction == proto.Direction_up:
		return proto.ErrTurnAbout
	case lastDirection == proto.Direction_left && *direction == proto.Direction_right:
		return proto.ErrTurnAbout
	case lastDirection == proto.Direction_right && *direction == proto.Direction_left:
		return proto.ErrTurnAbout
	}

	if len(snake.NextDirections) > buf {
		snake.NextDirections = append(snake.NextDirections[:buf], direction)
	} else {
		snake.NextDirections = append(snake.NextDirections, direction)
	}

	return nil
}

func generateDiffFromBody(s *proto.Snake) (diffs []*proto.Diff) {
	for _, bodyPart := range s.Body {
		diff := &proto.Diff{
			X:     bodyPart.X,
			Y:     bodyPart.Y,
			Color: s.Color,
			Add:   true,
		}

		diffs = append(diffs, diff)
	}

	return diffs
}
