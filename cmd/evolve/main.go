package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/taebow/evosnake/pkg/game"
	"github.com/taebow/evosnake/pkg/genetic"
	"github.com/taebow/evosnake/pkg/nn"
	"github.com/taebow/evosnake/pkg/nndriver"
	"github.com/taebow/evosnake/pkg/persist"
)

var nnConfig *nn.NeuralNetConfig = nn.NewNNConfig(8, 12, 12, 4)


func newPopulation(nnConfig *nn.NeuralNetConfig, size int) [][]float64 {
	pop := make([][]float64, size)
	for i := range pop {
		pop[i] = make([]float64, nnConfig.RawSize())
		for j := range pop[i] {
			pop[i][j] = (rand.Float64() * 2) - 1
		}
	}
	return pop
}

func PlayGame(rounds int, individual []float64) int {
	nnDriver := nndriver.NewNNDriver(nnConfig, individual)
	g := game.NewGame(50, 50, 45, 1, 1)
	g.Run(rounds, -1, false, nnDriver)
	maxScore := g.Snakes[0].MaxScore
	deaths := g.Snakes[0].Deaths
	return 100*maxScore - 10*(deaths*deaths)
}

func PlaySnakes(individuals [][]float64) {
	nnDrivers := make([]game.Driver, len(individuals))
	for i := range nnDrivers {
		nnDrivers[i] = nndriver.NewNNDriver(nnConfig, individuals[i])
	}
	// nnDriver := newNNDriver(nnConfig, individual)
	g := game.NewGame(50, 50, 5, len(individuals), 1)
	g.Run(-1, 25, true, nnDrivers...)
}

func MultiPlayGames(rounds int, individual []float64, nGames int) int {
	multiDriver := nndriver.NewMultiDriver([]*nn.NeuralNetConfig{nnConfig}, [][]float64{individual})
	games := make([]*game.Game, nGames)
	for i := range games {
		games[i] = game.NewGame(50, 50, 20, 1, 1)
	}
	game.RunMulti(games, rounds, multiDriver)
	fitnessSlice := make([]int, len(games))
	for i, g := range games {
		maxScore := g.Snakes[0].MaxScore
		deaths := g.Snakes[0].Deaths
		fitnessSlice[i] = 10*maxScore - (deaths*deaths)
	}
	return 10*min(fitnessSlice) + avg(fitnessSlice)
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

func avg(s []int) int {
	var sum int
	for _, v := range s {
		sum += v
	}
	return sum / len(s)
}

func std(s []int) int {
	s2 := make([]int, len(s))
	for i := range s2 {
		s2[i] = s[i]*s[i]
	}
	return int(math.Sqrt(float64(avg(s2) - avg(s))))
}


type FitFunc func (individual []float64) int

func train(nGenerations, genSize int, selectionRate, mutationRate float64, f FitFunc) ([][]float64, []int) {
	pop := newPopulation(nnConfig, genSize)
	var popBest, popMax [][]float64
	var fitBest, fitMax []int
	popFitness := make([]int, genSize)

	for nGen := 1; nGen <= nGenerations; nGen++ {
		var wg sync.WaitGroup
		for i, individual := range pop {
			wg.Add(1)
			go func(i int, individual []float64) {
				popFitness[i] = f(individual)
				wg.Done()
			}(i, individual)
		}
		wg.Wait()
		if len(fitBest) == 0 || max(popFitness) > min(fitBest) {
			popBest, fitBest = genetic.SelectRateBest(
				pop, 
				popFitness, 
				selectionRate,
			)
			popMax, fitMax = genetic.SelectNBest(
				append(popBest, popMax...),
				append(fitBest, fitMax...),
				len(popBest),
			)
		}
		_, fitMaxMax :=  genetic.SelectNBest(popMax, fitMax, 1)
		fmt.Printf("Trained generation %v, Max: %v, Elite %v\n", nGen, fitMaxMax, fitBest)
		popChild := genetic.Crossover(popBest, genSize)
		genetic.Mutate(popChild, mutationRate)
		pop = popChild
	}
	return popMax, fitMax
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	f := func (individual []float64) int {return MultiPlayGames(2000, individual, 5)}
	// f := func (individual []float64) int {return PlayGame(5000, individual)}
	records, fitRecords := train(2000, 100, 0.05, 0.1, f)
	record, _ := genetic.SelectNBest(records, fitRecords, 1)
	model := persist.NewModel(nnConfig.Layers, record[0])
	persist.Save("killer", model)
	PlaySnakes(record)
}
