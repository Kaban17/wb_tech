package config

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

// Config stores the configuration for the sort utility.
type Config struct {
	Column       int
	Numeric      bool
	Reverse      bool
	Unique       bool
	MonthSort    bool
	IgnoreTrail  bool
	CheckSorted  bool
	HumanNumeric bool
	InputFile    string
}

// New creates a new Config and parses command-line flags.
func New() *Config {
	cfg := &Config{}

	processedArgs := processCombinedFlags(os.Args[1:])

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	fs.IntVar(&cfg.Column, "k", 1, "sort by column")
	fs.BoolVar(&cfg.Numeric, "n", false, "sort by numeric value")
	fs.BoolVar(&cfg.Reverse, "r", false, "reverse the result of comparisons")
	fs.BoolVar(&cfg.Unique, "u", false, "output only the first of an equal run")
	fs.BoolVar(&cfg.MonthSort, "M", false, "compare month names")
	fs.BoolVar(&cfg.IgnoreTrail, "b", false, "ignore leading blanks")
	fs.BoolVar(&cfg.CheckSorted, "c", false, "check for sorted input")
	fs.BoolVar(&cfg.HumanNumeric, "h", false, "compare human readable numbers (e.g., 2K 1G)")

	fs.Parse(processedArgs)

	if fs.NArg() > 0 {
		cfg.InputFile = fs.Arg(0)
	}

	return cfg
}

func processCombinedFlags(args []string) []string {
	processedArgs := make([]string, 0, len(args))
	flagsWithValue := "k"

	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--") || len(arg) < 2 {
			processedArgs = append(processedArgs, arg)
			continue
		}

		// Pass through negative numbers
		if _, err := strconv.Atoi(arg); err == nil {
			processedArgs = append(processedArgs, arg)
			continue
		}

		chars := arg[1:]
		if len(chars) == 1 {
			processedArgs = append(processedArgs, arg)
			continue
		}

		for i, char := range chars {
			if strings.ContainsRune(flagsWithValue, char) {
				processedArgs = append(processedArgs, "-"+string(char))
				value := chars[i+1:]
				if len(value) > 0 {
					processedArgs = append(processedArgs, value)
				}
				goto nextArg
			} else {
				processedArgs = append(processedArgs, "-"+string(char))
			}
		}
	nextArg:
	}
	return processedArgs
}
