package lib

import "math"

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
