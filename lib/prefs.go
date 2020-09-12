package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type Person string

type Dish string

type Rating int

type Prefs map[Person]map[Dish]Rating

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
