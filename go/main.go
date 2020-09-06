package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

func main() {
	if err := Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}

func Run() error {
	args, err := ParseArgs()
	if err != nil {
		return fmt.Errorf("parse args: %w", err)
	}

	prefs, err := LoadPrefs(args.PrefsFileName)
	if err != nil {
		return fmt.Errorf("load preferences: %w", err)
	}
	index := BuildIndex(prefs)

	solutions := GenerateSolutions(index)
	for solution := range solutions {
		fmt.Println(solution)
	}

	return nil
}

type Args struct {
	PrefsFileName string
}

func ParseArgs() (*Args, error) {
	if len(os.Args) != 2 {
		return nil, fmt.Errorf("usage: %s PREFS", os.Args[0])
	}
	return &Args{PrefsFileName: os.Args[1]}, nil
}

type PersonName string
type DishName string
type Rating int

type Prefs map[PersonName]map[DishName]Rating

func LoadPrefs(fileName string) (Prefs, error) {
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var prefs Prefs
	if err := json.Unmarshal(fileBytes, &prefs); err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}
	if len(prefs) == 0 {
		return nil, errors.New("empty")
	}
	return prefs, nil
}

type Index struct {
	DishesByPerson Prefs
	PeopleByDish   map[DishName]map[PersonName]Rating
}

func (index Index) NumberOfPeople() int {
	return len(index.DishesByPerson)
}

func (index Index) NumberOfDishes() int {
	return len(index.PeopleByDish)
}

func BuildIndex(prefs Prefs) Index {
	peopleByDish := make(map[DishName]map[PersonName]Rating)
	for personName, dishes := range prefs {
		for dishName, rating := range dishes {
			people, ok := peopleByDish[dishName]
			if !ok {
				people = make(map[PersonName]Rating)
				peopleByDish[dishName] = people
			}
			people[personName] = rating
		}
	}
	return Index{
		DishesByPerson: prefs,
		PeopleByDish:   peopleByDish,
	}
}

type Solution string

func (sol Solution) Estimate() int {
	// TODO
	return 0
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
