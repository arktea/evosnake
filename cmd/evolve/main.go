package main

import (
	"math/rand"
	"time"

	"github.com/taebow/evosnake/pkg/genetic"
	"github.com/taebow/evosnake/pkg/nn"
	"github.com/taebow/evosnake/pkg/nndriver"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	modelConfig := nn.ModelConfig{8, 4}
	fitness := func (individual []float64) int {
		model := nn.NewModel(modelConfig, individual)
		games := nndriver.MultiPlayGames(500, model, 10)
		return genetic.EvaluateMultiGames(games)
	}
	solutions, fitSolutions := genetic.Train(1000, modelConfig.Size(), 100, 5, 0.1, fitness)
	best, _ := genetic.SelectBest(solutions, fitSolutions)
	model := nn.NewModel(modelConfig, best)
	nn.SaveModel("killer", model)
	nndriver.PlaySnakes([]*nn.Model{model})
}
