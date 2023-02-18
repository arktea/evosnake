package nndriver

import (
	"math"

	"github.com/taebow/evosnake/pkg/game"
)

func outputToDirection(output []float64) game.Direction {
	switch argMax(output) {
	case 0:
		return game.Up
	case 1:
		return game.Down
	case 2:
		return game.Left
	case 3:
		return game.Right
	}
	return game.Direction{}
}

func argMax(s []float64) int {
	max := -math.MaxFloat64
	var index int
	for i := range s {
		if s[i] > max {
			max = s[i]
			index = i
		}
	}
	return index
}