package reader

import (
	"bufio"
	"io"
	"os"

	"wb_tech/l2_10/internal/config"
)

// Read reads lines from the specified file or from STDIN.
func Read(cfg *config.Config) ([]string, error) {
	var reader io.Reader

	if cfg.InputFile != "" {
		file, err := os.Open(cfg.InputFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
