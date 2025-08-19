package types

import "strings"

type Command struct {
	Name           string
	Args           []string
	InputRedirect  string
	OutputRedirect string
	Append         bool
	PipeTo         *Command
	Background     bool
}

type Commands []*Command

func ToCommand(s string) Command {
	ss := strings.Split(s, " ")
	return Command{
		Name: ss[0],
		Args: ss[1:],
	}
}
