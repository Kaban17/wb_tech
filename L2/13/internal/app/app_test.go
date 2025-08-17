package app

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/pflag"
)

func TestRun(t *testing.T) {
	// Helper function to capture stdout
	captureOutput := func(f func()) string {
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		f()

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		io.Copy(&buf, r)
		return buf.String()
	}

	// Helper function to set command-line arguments
	setArgs := func(args ...string) {
		os.Args = append([]string{"cmd"}, args...)
	}

	testCases := []struct {
		name     string
		args     []string
		input    string
		expected string
	}{
		{
			name:     "Simple case",
			args:     []string{"-f", "1,3", "-d", ","},
			input:    "a,b,c\nd,e,f",
			expected: "a,c\nd,f\n",
		},
		{
			name:     "Range",
			args:     []string{"-f", "2-4", "-d", " "},
			input:    "a b c d e\nf g h i j",
			expected: "b c d\ng h i\n",
		},
		{
			name:     "Separated flag",
			args:     []string{"-f", "1,2", "-d", ",", "-s"},
			input:    "a,b\nc d\ne,f",
			expected: "a,b\ne,f\n",
		},
		{
			name:     "Default delimiter (tab)",
			args:     []string{"-f", "2"},
			input:    "a\tb\tc\nd\te\tf",
			expected: "b\ne\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
			// Set up stdin
			oldStdin := os.Stdin
			r, w, _ := os.Pipe()
			os.Stdin = r
			w.Write([]byte(tc.input))
			w.Close()

			// Set command-line arguments
			setArgs(tc.args...)

			// Capture output
			output := captureOutput(Run)

			// Restore stdin
			os.Stdin = oldStdin

			if output != tc.expected {
				t.Errorf("expected:\n%q\ngot:\n%q", tc.expected, output)
			}
		})
	}
}
