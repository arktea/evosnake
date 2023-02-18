package nndriver

import "github.com/taebow/evosnake/pkg/game"

func outputToDirection(output []float64) game.Direction {
	directionInt := argMax(output)
	var direction game.Direction
	switch directionInt {
	case 0:
		direction = game.Up
	case 1:
		direction = game.Down
	case 2:
		direction = game.Left
	case 3:
		direction = game.Right
	}
	return direction
}

func argMax(s []float64) int {
	var max float64
	var index int
	for i := range s {
		if s[i] > max {
			max = s[i]
			index = i
		}
	}
	return index
}