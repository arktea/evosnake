package nndriver

import (
	"github.com/taebow/evosnake/pkg/game"
	"github.com/taebow/evosnake/pkg/nn"
)

type MultiDriver struct {
	nn []*nn.NeuralNet
}

func NewMultiDriver(models []*nn.Model) *MultiDriver {
	neuralNets := make([]*nn.NeuralNet, len(models))
	for i, model := range models {
		neuralNets[i] = nn.NewNN(model)
	}
	return &MultiDriver{nn: neuralNets}
}

func (md *MultiDriver) GetDirections(games []*game.Game) [][]game.Direction {
	outputs := make([][][]float64, len(md.nn))
	directions := make([][]game.Direction, len(games))
	for i, nn := range md.nn {
		inputs := make([][]float64, len(games))
		for j, g := range games {
			inputs[j] =  gameToInput(g.Snakes[i], g.Foods[0], g.Board)
		}
		outputs[i] = nn.Predict(inputs...)
	}
	for i := range directions {
		directions[i] = make([]game.Direction, len(md.nn))
	}
	for i := range outputs {
		for j := range outputs[i] {
			directions[j][i] = outputToDirection(outputs[i][j])
		}
	}
	return directions
}