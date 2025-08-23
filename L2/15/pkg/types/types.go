package types

// Command represents a command in the shell.
type Command struct {
	Name           string
	Args           []string
	InputRedirect  string
	OutputRedirect string
	Append         bool
	PipeTo         *Command
	Background     bool
}
