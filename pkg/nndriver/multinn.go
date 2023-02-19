package nndriver

import (
	"github.com/taebow/evosnake/pkg/game"
	"github.com/taebow/evosnake/pkg/nn"
)

type MultiNNDriver struct {
	nn *nn.NeuralNet
}

func NewMultiNNDriver(model *nn.Model) *MultiNNDriver {
	return &MultiNNDriver{nn: nn.NewNN(model)}
}

func (md *MultiNNDriver) GetDirections(snakes []*game.Snake, games []*game.Game) []game.Direction {
	directions := make([]game.Direction, len(games))
	inputs := make([][]float64, len(games))
	for i, g := range games {
		inputs[i] =  gameToInput(snakes[i], g.Foods[0], g.Board)
	}
	outputs := md.nn.Predict(inputs...)
	for i := range outputs {
		directions[i] = outputToDirection(outputs[i])
	}
	return directions
}