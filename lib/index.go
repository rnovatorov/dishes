package lib

type Index struct {
	Menu   []DishName
	People []PersonName
	Matrix [][]Rating
}

func BuildIndex(prefs Prefs) Index {
	var index Index

	seen := make(map[DishName]bool)
	for personName, dishes := range prefs {
		for dishName := range dishes {
			if !seen[dishName] {
				index.Menu = append(index.Menu, dishName)
				seen[dishName] = true
			}
		}
		index.People = append(index.People, personName)
	}

	index.Matrix = make([][]Rating, len(index.People))
	for personIndex, personName := range index.People {
		index.Matrix[personIndex] = make([]Rating, len(index.Menu))
		for dishIndex, dishName := range index.Menu {
			index.Matrix[personIndex][dishIndex] = prefs[personName][dishName]
		}
	}

	return index
}
