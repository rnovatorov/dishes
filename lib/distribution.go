package lib

import (
	"fmt"
	"strconv"
	"strings"
)

type Distribution []int

func NewDistibution(index Index, i int) Distribution {
	iBaseNPeople := strconv.FormatInt(int64(i), len(index.People))
	nDishesZeroPadded := fmt.Sprintf("%%0%ds", len(index.Menu))
	digits := fmt.Sprintf(nDishesZeroPadded, iBaseNPeople)

	distribution := make([]int, len(index.Menu))
	for dishIndex, digit := range strings.Split(digits, "") {
		personIndex := int(MustParseInt(digit, len(index.People), 32))
		distribution[dishIndex] = personIndex
	}
	return distribution
}

func (d Distribution) Map(index Index) map[Person][]Dish {
	m := make(map[Person][]Dish)
	for dishIndex, personIndex := range d {
		person := index.People[personIndex]
		dish := index.Menu[dishIndex]
		m[person] = append(m[person], dish)
	}
	return m
}

func MustParseInt(s string, base int, bitSize int) int64 {
	n, err := strconv.ParseInt(s, base, bitSize)
	if err != nil {
		panic(err)
	}
	return n
}
