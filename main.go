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

	prefs, err := lib.LoadPrefs(*args.prefsFileName)
	if err != nil {
		return fmt.Errorf("load preferences: %w", err)
	}

	index := lib.BuildIndex(prefs)
	solutions := lib.GenerateSolutions(index)
	estimations := lib.EstimateSolutions(index, solutions)
	bestDistr, rating := lib.FindBestDistribution(estimations)

	fmt.Println(bestDistr.Map(index), rating)
	return nil
}

type parsedArgs struct {
	prefsFileName *string
	cpuProfile    *string
}

func parseArgs() parsedArgs {
	var args parsedArgs

	args.cpuProfile = flag.String("cpu-profile", "", "cpu profile file name")
	args.prefsFileName = flag.String("preferences", "", "preferences file name")

	flag.Parse()

	return args
}
