package lib

import (
	"math"
)

type Solution struct {
	Distribution Distribution
	Score        float64
}

func NewSolution(index Index, d Distribution) Solution {
	var cumRating float64
	personalRating := make([]float64, len(index.People))

	for dishIndex, personIndex := range d {
		rating := float64(index.Matrix[personIndex][dishIndex])
		cumRating += rating
		personalRating[personIndex] += rating
	}

	mean := cumRating / float64(len(index.People))

	var cumDeviation float64
	for _, rating := range personalRating {
		cumDeviation += math.Abs(mean - rating)
	}

	score := cumRating - cumDeviation

	return Solution{
		Distribution: d,
		Score:        score,
	}
}

func GenerateDistributions(index Index) <-chan Distribution {
	const bufSize = 256
	distributions := make(chan Distribution, bufSize)
	go func() {
		defer close(distributions)
		n := CountDistributions(index)
		for i := 0; i < n; i++ {
			distributions <- NewDistibution(index, i)
		}
	}()
	return distributions
}

func CountDistributions(index Index) int {
	nPeople := float64(len(index.People))
	nDishes := float64(len(index.Menu))
	n := math.Pow(nPeople, nDishes)
	return int(n)
}

func RateDistributions(index Index, distributions <-chan Distribution) <-chan Solution {
	const bufSize = 256
	solutions := make(chan Solution, bufSize)
	go func() {
		defer close(solutions)
		for d := range distributions {
			solutions <- NewSolution(index, d)
		}
	}()
	return solutions
}

func FindBestSolution(solutions <-chan Solution) Solution {
	best := Solution{Score: math.Inf(-1)}
	for s := range solutions {
		if s.Score > best.Score {
			best = s
		}
	}
	return best
}
