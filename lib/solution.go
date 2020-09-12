package lib

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Solution string

func (sol Solution) Score(index Index) float64 {
	var cumRating float64
	personalRating := make([]float64, len(index.People))

	sol.iterate(index, func(personIndex, dishIndex int) {
		rating := float64(index.Matrix[personIndex][dishIndex])
		cumRating += rating
		personalRating[personIndex] += rating
	})

	mean := cumRating / float64(len(index.People))

	var cumDeviation float64
	for _, rating := range personalRating {
		cumDeviation += math.Abs(mean - rating)
	}

	return cumRating - cumDeviation
}

type Distribution map[PersonName][]DishName

func (sol Solution) Distribution(index Index) Distribution {
	distr := make(Distribution)
	sol.iterate(index, func(personIndex, dishIndex int) {
		personName := index.People[personIndex]
		dishName := index.Menu[dishIndex]
		distr[personName] = append(distr[personName], dishName)
	})
	return distr
}

func (sol Solution) iterate(index Index, fn func(personIndex, dishIndex int)) {
	for dishIndex, s := range strings.Split(string(sol), "") {
		personIndex := int(MustParseInt(s, len(index.People), 32))
		fn(personIndex, dishIndex)
	}
}

func MustParseInt(s string, base int, bitSize int) int64 {
	n, err := strconv.ParseInt(s, base, bitSize)
	if err != nil {
		panic(err)
	}
	return n
}

func GenerateSolutions(index Index) <-chan Solution {
	const bufSize = 256
	solutions := make(chan Solution, bufSize)
	go func() {
		defer close(solutions)
		var i int64
		nSolutions := CountSolutions(index)
		for i = 0; i < nSolutions; i++ {
			iBaseNPeople := strconv.FormatInt(i, len(index.People))
			nDishesZeroPadded := fmt.Sprintf("%%0%ds", len(index.Menu))
			sol := fmt.Sprintf(nDishesZeroPadded, iBaseNPeople)
			solutions <- Solution(sol)
		}
	}()
	return solutions
}

func CountSolutions(index Index) int64 {
	nPeople := float64(len(index.People))
	nDishes := float64(len(index.Menu))
	n := math.Pow(nPeople, nDishes)
	return int64(n)
}

func FindBestDistribution(index Index, solutions <-chan Solution) (Distribution, float64) {
	var bestDistr Distribution
	maxScore := math.Inf(-1)
	for sol := range solutions {
		if score := sol.Score(index); score > maxScore {
			maxScore = score
			bestDistr = sol.Distribution(index)
		}
	}
	return bestDistr, maxScore
}
