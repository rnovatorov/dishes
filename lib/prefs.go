package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

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

type PersonNameSlice []PersonName

func (p PersonNameSlice) Len() int           { return len(p) }
func (p PersonNameSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p PersonNameSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type DishNameSlice []DishName

func (p DishNameSlice) Len() int           { return len(p) }
func (p DishNameSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p DishNameSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
