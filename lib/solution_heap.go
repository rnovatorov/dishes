package lib

import "container/heap"

type SolutionHeap struct {
	solutions []Solution
	capacity  int
}

func (h *SolutionHeap) Len() int {
	return len(h.solutions)
}

func (h *SolutionHeap) Less(i, j int) bool {
	return h.solutions[i].Score < h.solutions[j].Score
}

func (h *SolutionHeap) Swap(i, j int) {
	h.solutions[i], h.solutions[j] = h.solutions[j], h.solutions[i]
}

func (h *SolutionHeap) Push(x interface{}) {
	solution := x.(Solution)
	if h.capacity > 0 {
		switch n := h.Len(); n {
		case h.capacity:
			if min := h.solutions[0]; solution.Score > min.Score {
				h.solutions[0] = solution
				heap.Fix(h, 0)
			}
			return
		case 0:
			h.solutions = make([]Solution, 0, h.capacity)
		}
	}
	h.solutions = append(h.solutions, solution)
}

func (h *SolutionHeap) Pop() interface{} {
	n := h.Len()
	x := h.solutions[n-1]
	h.solutions = h.solutions[:n-1]
	return x
}
