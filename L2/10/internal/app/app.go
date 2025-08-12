package app

import (
	"fmt"
	"os"

	"wb_tech/l2_10/internal/config"
	"wb_tech/l2_10/internal/reader"
	"wb_tech/l2_10/internal/sorter"
)

// Run is the main application function.
func Run() {
	cfg := config.New()

	lines, err := reader.Read(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}

	sortedLines, err := sorter.Sort(lines, cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, line := range sortedLines {
		fmt.Println(line)
	}
}
