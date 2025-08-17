package app

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"wb_tech/l2_13/internal/config"
	"wb_tech/l2_13/internal/reader"
)

// Run starts the cut utility.
func Run() {
	cfg := config.New()
	scanner := reader.NewScanner()
	if err := Cut(scanner, cfg); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// Cut processes input from the scanner based on the provided config.
func Cut(scanner *bufio.Scanner, cfg *config.Config) error {
	fields, err := parseFields(cfg.Fields)
	if err != nil {
		return err
	}

	for scanner.Scan() {
		line := scanner.Text()

		if cfg.Separated && !strings.Contains(line, cfg.Delimiter) {
			continue
		}

		parts := strings.Split(line, cfg.Delimiter)
		var result []string

		for _, field := range fields {
			if field > 0 && field <= len(parts) {
				result = append(result, parts[field-1])
			}
		}

		if len(result) > 0 {
			fmt.Println(strings.Join(result, cfg.Delimiter))
		}
	}

	return scanner.Err()
}

// parseFields parses a fields string like "1,3-5" into a slice of integers.
func parseFields(fieldsStr string) ([]int, error) {
	if fieldsStr == "" {
		return nil, fmt.Errorf("fields cannot be empty")
	}

	var fields []int
	seen := make(map[int]bool)

	parts := strings.Split(fieldsStr, ",")
	for _, part := range parts {
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range: %s", part)
			}

			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid field number: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid field number: %s", rangeParts[1])
			}

			if start > end {
				return nil, fmt.Errorf("invalid range: start > end")
			}

			for i := start; i <= end; i++ {
				if !seen[i] {
					fields = append(fields, i)
					seen[i] = true
				}
			}
		} else {
			field, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid field number: %s", part)
			}
			if !seen[field] {
				fields = append(fields, field)
				seen[field] = true
			}
		}
	}
	return fields, nil
}
