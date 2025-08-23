package executor

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"wb_tech/l2_15/internal/parser"
)

func TestEcho(t *testing.T) {
	var buf bytes.Buffer
	echo([]string{"hello", "world"}, &buf)
	output := buf.String()

	if output != "hello world\n" {
		t.Errorf("echo failed: expected 'hello world\\n', got %q", output)
	}
}

func TestRedirects(t *testing.T) {
	// Test output redirect
	commands, _ := parser.Parse("echo hello > test.txt")
	Execute(commands)
	content, err := ioutil.ReadFile("test.txt")
	if err != nil {
		t.Errorf("Failed to read test.txt: %v", err)
	}
	if string(content) != "hello\n" {
		t.Errorf("Output redirect failed: expected 'hello\\n', got %q", string(content))
	}
	os.Remove("test.txt")

	// Test input redirect
	ioutil.WriteFile("test.txt", []byte("world"), 0644)
	commands, _ = parser.Parse("cat < test.txt")
	
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Execute(commands)

	w.Close()
	os.Stdout = old

	var buf [128]byte
	n, _ := r.Read(buf[:])
	output := string(buf[:n])

	if output != "world" {
		t.Errorf("Input redirect failed: expected 'world', got %q", output)
	}
	os.Remove("test.txt")
}

func TestPipeline(t *testing.T) {
	commands, _ := parser.Parse("echo 'hello world' | wc -w")
	
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Execute(commands)

	w.Close()
	os.Stdout = old

	var buf [128]byte
	n, _ := r.Read(buf[:])
	output := string(buf[:n])
	
	// Trim space because wc -w output can be weird
	if strings.TrimSpace(output) != "2" {
		t.Errorf("Pipeline failed: expected '2', got %q", output)
	}
}