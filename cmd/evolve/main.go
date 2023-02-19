package main

import (
	"github.com/taebow/evosnake/pkg/genetic"
	"github.com/taebow/evosnake/pkg/nn"
	"github.com/taebow/evosnake/pkg/nndriver"
)

func main() {
	modelConfig := nn.ModelConfig{8, 4}
	fitness := func (solution []float64) int {
		model := nn.NewModel(modelConfig, solution)
		games := nndriver.PlayOneSnakeMultiGames(5000, 10, model)
		return genetic.EvaluateMultiGames(games)
	}
	solutions, fitSolutions := genetic.Train(20, modelConfig.Size(), 100, 5, 0.05, fitness)
	best, _ := genetic.SelectBest(solutions, fitSolutions)
	model := nn.NewModel(modelConfig, best)
	nn.SaveModel("dumber", model)
	nndriver.PlaySnakes(-1, model)
}
