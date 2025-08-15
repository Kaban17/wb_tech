package app

import (
	"fmt"
	"regexp"
	"wb_tech/l2_12/internal/config"
	"wb_tech/l2_12/internal/reader"
)

func Run() {
	cfg := config.New()
	lines, err := reader.Read(cfg)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	grep(lines, cfg)
}

func grep(lines []string, cfg *config.Config) {
	pattern := cfg.Regex
	if cfg.IgnoreCase {
		pattern = "(?i)" + pattern
	}
	if cfg.FixString {
		pattern = regexp.QuoteMeta(pattern)
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}

	// Handle -C flag
	if cfg.C > 0 {
		cfg.A = cfg.C
		cfg.B = cfg.C
	}

	// Find indices of all matching lines
	matchingIndices := []int{}
	for i, line := range lines {
		match := re.MatchString(line)
		if (match && !cfg.Invert) || (!match && cfg.Invert) {
			matchingIndices = append(matchingIndices, i)
		}
	}

	// If -c flag is set, print count and exit
	if cfg.CountString {
		fmt.Println(len(matchingIndices))
		return
	}

	// Build the set of lines to print, including context
	linesToPrint := make(map[int]bool)
	for _, idx := range matchingIndices {
		// Add the matching line itself
		linesToPrint[idx] = true

		// Context Before (-B)
		for i := 1; i <= cfg.B; i++ {
			if idx-i >= 0 {
				linesToPrint[idx-i] = true
			}
		}
		// Context After (-A)
		for i := 1; i <= cfg.A; i++ {
			if idx+i < len(lines) {
				linesToPrint[idx+i] = true
			}
		}
	}

	// Print the selected lines
	for i := 0; i < len(lines); i++ {
		if linesToPrint[i] {
			if cfg.LineNum {
				fmt.Printf("%d:", i+1)
			}
			fmt.Println(lines[i])
		}
	}
}
