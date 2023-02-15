package main

import (
	// "math"
	"math"
	"math/rand"
	"time"
	"fmt"
	
	"gonum.org/v1/gonum/mat"

	"github.com/taebow/evosnake/pkg/game"
	"github.com/taebow/evosnake/pkg/genetic"
	"github.com/taebow/evosnake/pkg/nn"
)

type NNDriver struct {
	nn *nn.NeuralNet
}

func newNNDriver(nnConfig *nn.NeuralNetConfig, weights []float64) *NNDriver {
	nn := nn.NewNN(nnConfig)
	nn.InitFromSlice(weights)
	return &NNDriver{nn: nn}
}

func (d *NNDriver) GetDirection(s *game.Snake, g *game.Game) game.Direction {
	inputs := s.See(g.Foods[0], g.Board)
	inputsMat := mat.NewDense(1, len(inputs), inputs)
	outputs, _ := d.nn.Predict(inputsMat)
	directionInt := argMax(outputs.RawMatrix().Data)
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

func newPopulation(genes, size int) [][]float64 {
	pop := make([][]float64, size)
	for i := range pop {
		pop[i] = make([]float64, genes)
		for j := range pop[i] {
			pop[i][j] = (rand.Float64() * 2) - 1
		}
	}
	return pop
}

func PlayGame(rounds int, individual []float64) int {
	nnConfig := nn.NewNNConfig(8, 4, 50)
	nnDriver := newNNDriver(nnConfig, individual)
	g := game.NewGame(50, 50, 20, 1, 1)
	g.Run(rounds, -1, false, []game.Driver{nnDriver})
	return g.Snakes[0].MaxScore * 100 - 10*(g.Snakes[0].Deaths*g.Snakes[0].Deaths)
}

func PlaySnake(individual []float64) {
	nnConfig := nn.NewNNConfig(8, 4, 50)
	nnDriver := newNNDriver(nnConfig, individual)
	g := game.NewGame(50, 50, 5, 1, 1)
	g.Run(-1, 25, true, []game.Driver{nnDriver})
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

func max(s []int) int {
	m := math.MinInt
	for _, v := range s {
		if v > m {
			m = v
		}
	}
	return m
}

func min(s []int) int {
	m := math.MaxInt
	for _, v := range s {
		if v < m {
			m = v
		}
	}
	return m
}

func train(nGenerations int) []float64 {
	pop := newPopulation(9*50+51*4, 100)
	var popBest [][]float64
	var popFitness, fitBest []int
	// var record []float64
	var fitnessRecord int = math.MinInt
	for nGen := 1; nGen <= nGenerations; nGen++ {
		popFitness = make([]int, 100)
		for i, individual := range pop {
			popFitness[i] = PlayGame(1000, individual)
		}
		if max(popFitness) > min(fitBest) || len(fitBest) == 0 {
			popBest, fitBest = genetic.SelectBest(pop, popFitness, 10)
		}
		if _, f := genetic.SelectBest(popBest, fitBest, 1); f[0] > fitnessRecord {
			// record = r[0]
			fitnessRecord = f[0]
		} 
		fmt.Printf("Trained generation %v, Record: %v\n", nGen, fitnessRecord)
		popChild := genetic.Crossover(popBest, 90)
		genetic.Mutate(popChild, 20)
		pop = append(popBest, popChild...)
	}
	popBest, _ = genetic.SelectBest(popBest, fitBest, 1)
	return popBest[0]
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	best := train(5000)
	PlaySnake(best)
}
