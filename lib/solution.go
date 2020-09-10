package lib

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Solution string

type Distribution [][]int

func (d Distribution) Map(index Index) map[PersonName][]DishName {
	m := make(map[PersonName][]DishName)
	for personIndex, dishes := range d {
		personName := index.People[personIndex]
		for _, dishIndex := range dishes {
			dishName := index.Menu[dishIndex]
			m[personName] = append(m[personName], dishName)
		}
	}
	return m
}

type Estimation struct {
	Rating       Rating
	Distribution Distribution
}

func (sol Solution) Estimate(index Index) Estimation {
	var rating Rating

	distr := make(Distribution, len(index.People))
	for personIndex := range index.People {
		distr[personIndex] = make([]int, 0, len(index.Menu))
	}

	for dishIndex, s := range strings.Split(string(sol), "") {
		personIndex := int(MustParseInt(s, len(index.People), 32))
		rating += index.Matrix[personIndex][dishIndex]
		distr[personIndex] = append(distr[personIndex], dishIndex)
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

		var i int64
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
