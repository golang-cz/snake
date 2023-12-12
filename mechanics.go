package main

import (
	"github.com/golang-cz/snake/proto"
)

func turnSnake(snake *proto.Snake, direction *proto.Direction, buf int) error {
	lastDirection := *snake.Direction
	if len(snake.NextDirections) > 0 {
		lastDirection = *snake.NextDirections[len(snake.NextDirections)-1]
	}

	// Same direction.
	if lastDirection == *direction {
		return proto.ErrInvalidTurn
	}

	// Disallow turnabouts.
	switch {
	case lastDirection == proto.Direction_up && *direction == proto.Direction_down:
		return proto.ErrInvalidTurn
	case lastDirection == proto.Direction_down && *direction == proto.Direction_up:
		return proto.ErrInvalidTurn
	case lastDirection == proto.Direction_left && *direction == proto.Direction_right:
		return proto.ErrInvalidTurn
	case lastDirection == proto.Direction_right && *direction == proto.Direction_left:
		return proto.ErrInvalidTurn
	}

	if len(snake.NextDirections) > buf {
		snake.NextDirections = append(snake.NextDirections[:buf], direction)
	} else {
		snake.NextDirections = append(snake.NextDirections, direction)
	}

	return nil
}
