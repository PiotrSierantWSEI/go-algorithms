package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

const (
	populationSize      = 100
	generations         = 300
	offspringPairs      = populationSize / 2
	tournamentSize      = 3
	mutationProbability = 0.35
	mutationStep        = 0.001
	period              = 2 * math.Pi
)

type individual struct {
	x       float64
	y       float64
	fitness float64
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	population := createInitialPopulation(rng)
	best := bestIndividual(population)

	fmt.Println("Algorytm ewolucyjny - maksimum funkcji")
	fmt.Println("f(x,y) = |sin(x)+sin(2x)+sin(4x)+sin(8x)| + |cos(y)+cos(2y)+cos(4y)+cos(8y)|")
	fmt.Printf("Populacja: %d, pokolenia: %d, mutacja: +/- %.3f\n", populationSize, generations, mutationStep)
	fmt.Println()

	for generation := 1; generation <= generations; generation++ {
		children := make([]individual, 0, populationSize)

		for pair := 0; pair < offspringPairs; pair++ {
			parentA := tournamentSelect(population, rng)
			parentB := tournamentSelect(population, rng)
			childA, childB := crossover(parentA, parentB)

			mutate(&childA, rng)
			mutate(&childB, rng)
			evaluate(&childA)
			evaluate(&childB)

			children = append(children, childA, childB)
		}

		candidates := make([]individual, 0, len(population)+len(children))
		candidates = append(candidates, population...)
		candidates = append(candidates, children...)
		population = selectNextGeneration(candidates, populationSize, rng)

		currentBest := bestIndividual(population)
		if currentBest.fitness > best.fitness {
			best = currentBest
		}

		if generation == 1 || generation%50 == 0 || generation == generations {
			fmt.Printf("Pokolenie %3d: x = %.6f, y = %.6f, f(x,y) = %.6f\n",
				generation, best.x, best.y, best.fitness)
		}
	}

	svgPath, pngPath, err := generateResultImages(best)
	if err != nil {
		fmt.Printf("Nie udalo sie wygenerowac wykresow: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("Najlepsze znalezione rozwiazanie:")
	fmt.Printf("x = %.9f\n", best.x)
	fmt.Printf("y = %.9f\n", best.y)
	fmt.Printf("f(x,y) = %.9f\n", best.fitness)
	fmt.Printf("Wykres SVG zapisano do: %s\n", svgPath)
	fmt.Printf("Wykres PNG dla raportu zapisano do: %s\n", pngPath)
}

func createInitialPopulation(rng *rand.Rand) []individual {
	population := make([]individual, populationSize)
	for i := range population {
		population[i] = individual{
			x: rng.Float64() * period,
			y: rng.Float64() * period,
		}
		evaluate(&population[i])
	}
	return population
}

func objective(x, y float64) float64 {
	sinPart := math.Sin(x) + math.Sin(2*x) + math.Sin(4*x) + math.Sin(8*x)
	cosPart := math.Cos(y) + math.Cos(2*y) + math.Cos(4*y) + math.Cos(8*y)
	return math.Abs(sinPart) + math.Abs(cosPart)
}

func evaluate(candidate *individual) {
	candidate.x = normalize(candidate.x)
	candidate.y = normalize(candidate.y)
	candidate.fitness = objective(candidate.x, candidate.y)
}

func crossover(parentA, parentB individual) (individual, individual) {
	childA := individual{
		x: parentA.x,
		y: (parentA.y + parentB.y) / 2,
	}

	childB := individual{
		x: (parentA.x + parentB.x) / 2,
		y: parentB.y,
	}

	return childA, childB
}

func mutate(candidate *individual, rng *rand.Rand) {
	if rng.Float64() < mutationProbability {
		candidate.x += randomSignedStep(rng)
	}
	if rng.Float64() < mutationProbability {
		candidate.y += randomSignedStep(rng)
	}
}

func randomSignedStep(rng *rand.Rand) float64 {
	if rng.Intn(2) == 0 {
		return -mutationStep
	}
	return mutationStep
}

func selectNextGeneration(candidates []individual, size int, rng *rand.Rand) []individual {
	sortedCandidates := append([]individual(nil), candidates...)
	sort.Slice(sortedCandidates, func(i, j int) bool {
		return sortedCandidates[i].fitness > sortedCandidates[j].fitness
	})

	nextGeneration := make([]individual, 0, size)
	eliteCount := min(2, len(sortedCandidates))
	nextGeneration = append(nextGeneration, sortedCandidates[:eliteCount]...)

	for len(nextGeneration) < size {
		nextGeneration = append(nextGeneration, tournamentSelect(candidates, rng))
	}

	return nextGeneration
}

func tournamentSelect(population []individual, rng *rand.Rand) individual {
	best := population[rng.Intn(len(population))]
	for i := 1; i < tournamentSize; i++ {
		candidate := population[rng.Intn(len(population))]
		if candidate.fitness > best.fitness {
			best = candidate
		}
	}
	return best
}

func bestIndividual(population []individual) individual {
	best := population[0]
	for _, candidate := range population[1:] {
		if candidate.fitness > best.fitness {
			best = candidate
		}
	}
	return best
}

func normalize(value float64) float64 {
	value = math.Mod(value, period)
	if value < 0 {
		value += period
	}
	return value
}
