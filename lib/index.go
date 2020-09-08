package lib

import (
	"sort"
)

type Index struct {
	DishesByPerson Prefs
	PeopleByDish   map[DishName]map[PersonName]Rating
	SortedDishes   []DishName
	SortedPeople   []PersonName
}

func (index Index) NumberOfPeople() int {
	return len(index.DishesByPerson)
}

func (index Index) NumberOfDishes() int {
	return len(index.PeopleByDish)
}

func BuildIndex(prefs Prefs) Index {
	dishesByPerson := prefs
	peopleByDish := make(map[DishName]map[PersonName]Rating)

	for personName, dishes := range dishesByPerson {
		for dishName, rating := range dishes {
			people, ok := peopleByDish[dishName]
			if !ok {
				people = make(map[PersonName]Rating)
				peopleByDish[dishName] = people
			}
			people[personName] = rating
		}
	}

	sortedDishes := make(DishNameSlice, 0, len(peopleByDish))
	for dishName := range peopleByDish {
		sortedDishes = append(sortedDishes, dishName)
	}
	sort.Sort(sortedDishes)

	sortedPeople := make(PersonNameSlice, 0, len(dishesByPerson))
	for personName := range dishesByPerson {
		sortedPeople = append(sortedPeople, personName)
	}
	sort.Sort(sortedPeople)

	return Index{
		DishesByPerson: dishesByPerson,
		PeopleByDish:   peopleByDish,
		SortedDishes:   sortedDishes,
		SortedPeople:   sortedPeople,
	}
}
