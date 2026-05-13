package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func f(x float64) float64 {
	return math.Abs(math.Sin(x) + math.Sin(2*x) + math.Sin(4*x) + math.Sin(8*x))
}

func main() {
	var epsilon float64
	fmt.Print("Podaj epsilon: ")
	fmt.Scan(&epsilon)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Prostokąt: x ∈ (0, 2π), y ∈ (0, 4), pole = 8π
	rectangleArea := 8 * math.Pi

	var L1, L2 int64 // L1: pod krzywą, L2: nad krzywą
	var prevIntegral float64
	var currentIntegral float64
	const minPoints = 1000

	for {
		x := rng.Float64() * 2 * math.Pi
		y := rng.Float64() * 4

		if y <= f(x) {
			L1++
		} else {
			L2++
		}

		total := L1 + L2
		currentIntegral = rectangleArea * float64(L1) / float64(total)

		if total >= minPoints && math.Abs(currentIntegral-prevIntegral) < epsilon {
			break
		}

		prevIntegral = currentIntegral
	}

	fmt.Printf("Wynik całki: %.6f\n", currentIntegral)
	fmt.Printf("Liczba punktów: %d\n", L1+L2)
}
