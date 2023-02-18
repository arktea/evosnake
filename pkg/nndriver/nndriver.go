package nndriver

import (
	"github.com/taebow/evosnake/pkg/game"
	"github.com/taebow/evosnake/pkg/nn"
)

type NNDriver struct {
	nn *nn.NeuralNet
}

func NewNNDriver(nnConfig *nn.NeuralNetConfig, weights []float64) *NNDriver {
	nn := nn.NewNN(nnConfig, weights)
	return &NNDriver{nn: nn}
}

func (d *NNDriver) GetDirection(s *game.Snake, g *game.Game) game.Direction {
	inputs := gameToInput(s, g.Foods[0], g.Board)
	outputs := d.nn.Predict(inputs)
	return outputToDirection(outputs[0])
}