package lib

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
