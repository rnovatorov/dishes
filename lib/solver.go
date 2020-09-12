package lib

import (
	"container/heap"
	"math"
)

type Solver struct {
	index Index
}

func NewSolver(index Index) Solver {
	return Solver{index: index}
}

func (s Solver) CountDistributions() int {
	nPeople := float64(len(s.index.People))
	nDishes := float64(len(s.index.Menu))
	n := math.Pow(nPeople, nDishes)
	return int(n)
}

func (s Solver) GenerateDistributions() <-chan Distribution {
	const bufSize = 256
	distributions := make(chan Distribution, bufSize)
	go func() {
		defer close(distributions)
		n := s.CountDistributions()
		for i := 0; i < n; i++ {
			distributions <- NewDistibution(s.index, i)
		}
	}()
	return distributions
}

func (s Solver) RateDistributions(distributions <-chan Distribution) <-chan Solution {
	const bufSize = 256
	solutions := make(chan Solution, bufSize)
	go func() {
		defer close(solutions)
		for d := range distributions {
			solutions <- NewSolution(s.index, d)
		}
	}()
	return solutions
}

func (s Solver) FindBestSolutions(solutions <-chan Solution, n int) []Solution {
	h := &SolutionHeap{capacity: n}
	for sol := range solutions {
		heap.Push(h, sol)
	}
	best := make([]Solution, n)
	for i := n - 1; i >= 0; i-- {
		best[i] = heap.Pop(h).(Solution)
	}
	return best
}
