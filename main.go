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

	if *args.cpuProf != "" {
		file, err := os.Create(*args.cpuProf)
		if err != nil {
			return fmt.Errorf("create file: %w", err)
		}
		defer file.Close()
		if err := pprof.StartCPUProfile(file); err != nil {
			return fmt.Errorf("start pprof: %w", err)
		}
		defer pprof.StopCPUProfile()
	}

	prefs, err := lib.LoadPreferences(*args.prefs)
	if err != nil {
		return fmt.Errorf("load preferences: %w", err)
	}

	index := lib.BuildIndex(prefs)
	if *args.norm {
		index.Normalize()
	}
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
	prefs   *string
	cpuProf *string
	top     *uint
	norm    *bool
}

func parseArgs() parsedArgs {
	var args parsedArgs

	args.cpuProf = flag.String("cpu-prof", "", "cpu profile file name")
	args.prefs = flag.String("prefs", "", "preferences file name")
	args.top = flag.Uint("top", 10, "find top n solutions")
	args.norm = flag.Bool("norm", false, "normalize preferences weights")

	flag.Parse()

	return args
}
