package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type Person string

type Dish string

type Rating float64

type Preferences map[Person]map[Dish]Rating

func LoadPreferences(fileName string) (Preferences, error) {
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	var prefs Preferences
	if err := json.Unmarshal(fileBytes, &prefs); err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}
	if len(prefs) == 0 {
		return nil, errors.New("empty")
	}
	return prefs, nil
}
