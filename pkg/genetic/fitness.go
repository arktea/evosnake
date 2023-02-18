package genetic

import (
	"sync"

	"github.com/taebow/evosnake/pkg/game"
)


type FitnessFunc func (solution []float64) int

func evalFitnessParallel(pop [][]float64, fit FitnessFunc) []int {
	var wg sync.WaitGroup
	popFitness := make([]int, len(pop))
	for i, solution := range pop {
		wg.Add(1)
		go func(i int, solution []float64) {
			popFitness[i] = fit(solution)
			wg.Done()
		}(i, solution)
	}
	wg.Wait()
	return popFitness
}

func evaluateGame(g *game.Game) int {
	maxScore := g.Snakes[0].MaxScore
	deaths := g.Snakes[0].Deaths
	return 10*maxScore - deaths*deaths
}

func EvaluateMultiGames(games []*game.Game) int {
	fitness := make([]int, len(games))
	for i, g := range games {
		fitness[i] = evaluateGame(g)
	}
	return 10*Min(fitness) + Avg(fitness)
}