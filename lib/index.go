package lib

import "math"

type Index struct {
	Menu   []Dish
	People []Person
	Matrix [][]Rating
}

func BuildIndex(prefs Preferences) Index {
	var index Index

	seen := make(map[Dish]bool)
	for person, dishes := range prefs {
		for dish := range dishes {
			if !seen[dish] {
				index.Menu = append(index.Menu, dish)
				seen[dish] = true
			}
		}
		index.People = append(index.People, person)
	}

	index.Matrix = make([][]Rating, len(index.People))
	for personIndex, person := range index.People {
		index.Matrix[personIndex] = make([]Rating, len(index.Menu))
		for dishIndex, dish := range index.Menu {
			index.Matrix[personIndex][dishIndex] = prefs[person][dish]
		}
	}

	return index
}

func (index Index) Normalize() {
	for _, dishes := range index.Matrix {
		minRating := Rating(math.Inf(1))
		cumRating := Rating(0)
		for _, rating := range dishes {
			if rating < minRating {
				minRating = rating
			}
			cumRating += rating
		}

		shift := Rating(0)
		switch {
		case minRating < 0:
			shift = -minRating
		case minRating == 0 && cumRating == 0:
			shift = 1
		}
		if shift != 0 {
			for dishIndex := range dishes {
				dishes[dishIndex] += shift
				cumRating += shift
			}
		}

		for dishIndex := range dishes {
			dishes[dishIndex] /= cumRating
		}
	}
}
