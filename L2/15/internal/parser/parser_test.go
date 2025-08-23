package parser

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input              string
		expectedCmds       int
		expectedInput      string
		expectedOutput     string
		expectedAppend     bool
	}{
		{"ls -l", 1, "", "", false},
		{"ps aux | grep myprocess", 2, "", "", false},
		{"echo 'hello' > file.txt", 1, "", "file.txt", false},
		{"cat < file.txt", 1, "file.txt", "", false},
		{"cat < file.txt | grep hello > out.txt", 2, "file.txt", "out.txt", false},
	}

	for _, test := range tests {
		commands, err := Parse(test.input)
		if err != nil {
			t.Errorf("Parse(%q) returned error: %v", test.input, err)
		}
		if len(commands) != test.expectedCmds {
			t.Errorf("Parse(%q) = %d commands, want %d", test.input, len(commands), test.expectedCmds)
		}
		if len(commands) > 0 {
			if commands[0].InputRedirect != test.expectedInput {
				t.Errorf("Parse(%q) input = %q, want %q", test.input, commands[0].InputRedirect, test.expectedInput)
			}
			lastCmd := commands[len(commands)-1]
			if lastCmd.OutputRedirect != test.expectedOutput {
				t.Errorf("Parse(%q) output = %q, want %q", test.input, lastCmd.OutputRedirect, test.expectedOutput)
			}
			if lastCmd.Append != test.expectedAppend {
				t.Errorf("Parse(%q) append = %v, want %v", test.input, lastCmd.Append, test.expectedAppend)
			}
		}
	}
}