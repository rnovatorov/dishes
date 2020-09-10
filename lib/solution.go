package lib

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Solution string

func (sol Solution) Rating(index Index) Rating {
	var rating Rating
	sol.iterate(index, func(personIndex, dishIndex int) {
		rating += index.Matrix[personIndex][dishIndex]
	})
	return rating
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

func FindBestDistribution(index Index, solutions <-chan Solution) (Distribution, Rating) {
	var bestDistr Distribution
	maxRating := Rating(math.Inf(-1))
	for sol := range solutions {
		if rating := sol.Rating(index); rating > maxRating {
			maxRating = rating
			bestDistr = sol.Distribution(index)
		}
	}
	return bestDistr, maxRating
}
