package genetic

import (
	"fmt"
	"time"
	"math/rand"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func generateRandPopulation(solSize, popSize int) [][]float64 {
	pop := make([][]float64, popSize)
	for i := range pop {
		pop[i] = make([]float64, solSize)
		for j := range pop[i] {
			pop[i][j] = (rand.Float64() * 2) - 1
		}
	}
	return pop
}

func Train(nGen, solSize, popSize, selectSize int, mutationRate float64, f FitnessFunc) ([][]float64, []int) {
	var popBest, popElite [][]float64
	var fitBest, fitElite []int
	pop := generateRandPopulation(solSize, popSize)
	for gen := 1; gen <= nGen; gen++ {
		popFitness := evalFitnessParallel(pop, f)
		if len(fitBest) == 0 || Max(popFitness) > Min(fitBest) {
			popBest, fitBest = selectNBest(pop, popFitness, selectSize)
			popElite, fitElite = selectNBest(
				append(popBest, popElite...),
				append(fitBest, fitElite...),
				len(popBest),
			)
		}
		_, fitMaxElite :=  SelectBest(popElite, fitElite)
		popChild := crossover(popBest, popSize)
		mutate(popChild, mutationRate)
		pop = popChild
		fmt.Printf("Trained generation %v, Max: %v, Elite: %v, Best: %v\n", gen, fitMaxElite, fitElite, fitBest)
	}
	return popElite, fitElite
}