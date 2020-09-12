package lib

type Distribution []int

func NewDistibution(index Index, i int) Distribution {
	digits := changeNumeralSystem(i, len(index.People), len(index.Menu))
	distribution := make([]int, len(index.Menu))
	for dishIndex, personIndex := range digits {
		distribution[dishIndex] = personIndex
	}
	return distribution
}

func changeNumeralSystem(n int, base int, size int) []int {
	digits := make([]int, size)
	for i := size - 1; n != 0; {
		q, r := n/base, n%base
		digits[i] = r
		i--
		n = q
	}
	return digits
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
