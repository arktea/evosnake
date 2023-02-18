package genetic

import (
	"math"
	"math/rand"
)

func selectNBest(pop [][]float64, popFitness []int, n int) ([][]float64, []int) {
	bestFitness, bestSolution := make([]int, n), make([][]float64, n)
	for i := range bestFitness {
		bestFitness[i] = math.MinInt
	}
	for i, fitness := range popFitness {
		minBest := math.MaxInt
		var minBestIndex int
		for j, best := range bestFitness {
			if best < minBest {
				minBestIndex = j
				minBest = best
			}
		}
		if fitness > bestFitness[minBestIndex] {
			bestFitness[minBestIndex] = fitness
			bestSolution[minBestIndex] = pop[i]
		}
	}
	return bestSolution, bestFitness
}

func SelectBest(pop [][]float64, popFitness []int) ([]float64, int) {
	best, fitness := selectNBest(pop, popFitness, 1)
	return best[0], fitness[0]
}


func crossover(pop [][]float64, n int) [][]float64 {
	res := make([][]float64, n)
	for i := range res {
		parent1 := pop[rand.Intn(len(pop))]
		parent2 := pop[rand.Intn(len(pop))]
		child := make([]float64, len(parent1))
		for j := range child {
			if rand.Intn(2) == 1 {
				child[j] = parent1[j]
			} else {
				child[j] = parent2[j]
			}
		}
		res[i] = child
	}
	return res
}

func mutate(pop [][]float64, rate float64) {
	n := int(float64(len(pop)) * rate)
	for i := range pop {
		for j := 0; j < rand.Intn(n); j++ {
			pop[i][rand.Intn(len(pop[i]))] = (rand.Float64() * 2) - 1
		}
	}
}
