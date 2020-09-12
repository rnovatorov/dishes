package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/rnovatorov/dishes/lib"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	args := parseArgs()

	if *args.top == 0 {
		return nil
	}

	if *args.cpuProfile != "" {
		file, err := os.Create(*args.cpuProfile)
		if err != nil {
			return fmt.Errorf("create file: %w", err)
		}
		defer file.Close()
		if err := pprof.StartCPUProfile(file); err != nil {
			return fmt.Errorf("start pprof: %w", err)
		}
		defer pprof.StopCPUProfile()
	}

	prefs, err := lib.LoadPreferences(*args.prefsFileName)
	if err != nil {
		return fmt.Errorf("load preferences: %w", err)
	}

	index := lib.BuildIndex(prefs)
	solver := lib.NewSolver(index)

	distributions := solver.GenerateDistributions()
	solutions := solver.RateDistributions(distributions)
	bestSolutions := solver.FindBestSolutions(solutions, int(*args.top))

	for i, s := range bestSolutions {
		fmt.Printf("%02d %v %f\n", i+1, s.Distribution.Map(index), s.Score)
	}

	return nil
}

type parsedArgs struct {
	prefsFileName *string
	cpuProfile    *string
	top           *uint
}

func parseArgs() parsedArgs {
	var args parsedArgs

	args.cpuProfile = flag.String("cpu-profile", "", "cpu profile file name")
	args.prefsFileName = flag.String("preferences", "", "preferences file name")
	args.top = flag.Uint("top", 10, "find top n solutions")

	flag.Parse()

	return args
}
