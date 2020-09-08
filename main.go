package main

import (
	"fmt"
	"os"

	"github.com/rnovatorov/dishes/lib"
)

func main() {
	if err := Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func Run() error {
	args, err := ParseArgs()
	if err != nil {
		return fmt.Errorf("parse args: %w", err)
	}

	prefs, err := lib.LoadPrefs(args.PrefsFileName)
	if err != nil {
		return fmt.Errorf("load preferences: %w", err)
	}
	index := lib.BuildIndex(prefs)

	solutions := lib.GenerateSolutions(index)
	estimations := lib.EstimateSolutions(index, solutions)
	bestDistr, rating := lib.FindBestDistribution(estimations)

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
