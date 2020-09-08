package lib

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Solution string

type Distribution map[PersonName][]DishName

type Estimation struct {
	Rating       Rating
	Distribution Distribution
}

func (sol Solution) Estimate(index Index) Estimation {
	var rating Rating
	distr := make(map[PersonName][]DishName)

	for dishIndex, s := range strings.Split(string(sol), "") {
		personIndex := MustParseInt(s, index.NumberOfPeople(), 32)

		dishName := index.SortedDishes[dishIndex]
		personName := index.SortedPeople[personIndex]

		rating += index.DishesByPerson[personName][dishName]
		distr[personName] = append(distr[personName], dishName)
	}

	return Estimation{Rating: rating, Distribution: distr}
}

func MustParseInt(s string, base int, bitSize int) int64 {
	n, err := strconv.ParseInt(s, base, bitSize)
	if err != nil {
		panic(err)
	}
	return n
}

func GenerateSolutions(index Index) <-chan Solution {
	solutions := make(chan Solution)
	go func() {
		defer close(solutions)

		nSolutions := CountSolutions(index)
		nPeople := index.NumberOfPeople()
		nDishes := index.NumberOfDishes()

		var i int64
		for i = 0; i < nSolutions; i++ {
			iBaseNPeople := strconv.FormatInt(i, nPeople)
			nDishesZeroPadded := fmt.Sprintf("%%0%ds", nDishes)
			sol := fmt.Sprintf(nDishesZeroPadded, iBaseNPeople)
			solutions <- Solution(sol)
		}
	}()
	return solutions
}

func CountSolutions(index Index) int64 {
	nPeople := float64(index.NumberOfPeople())
	nDishes := float64(index.NumberOfDishes())
	n := math.Pow(nPeople, nDishes)
	return int64(n)
}

func EstimateSolutions(index Index, solutions <-chan Solution) <-chan Estimation {
	estimations := make(chan Estimation)
	go func() {
		defer close(estimations)
		for solution := range solutions {
			estimations <- solution.Estimate(index)
		}
	}()
	return estimations
}

func FindBestDistribution(estimations <-chan Estimation) (Distribution, Rating) {
	best := Estimation{Rating: Rating(math.Inf(-1))}
	for estimation := range estimations {
		if estimation.Rating > best.Rating {
			best = estimation
		}
	}
	return best.Distribution, best.Rating
}
