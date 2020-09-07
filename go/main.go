package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
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
	estimations := EstimateSolutions(index, solutions)
	bestDistr, rating := FindBestDistribution(estimations)

	fmt.Println(bestDistr, rating)
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

type PersonNameSlice []PersonName

func (p PersonNameSlice) Len() int           { return len(p) }
func (p PersonNameSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p PersonNameSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type DishNameSlice []DishName

func (p DishNameSlice) Len() int           { return len(p) }
func (p DishNameSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p DishNameSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
