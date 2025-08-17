package reader

import (
	"bufio"
	"os"
)

// NewScanner creates a new scanner for reading from STDIN.
func NewScanner() *bufio.Scanner {
	return bufio.NewScanner(os.Stdin)
}
