package parser

import (
	"strings"
	"wb_tech/l2_15/pkg/types"
)

// Parse parses a string into a slice of commands.
func Parse(input string) ([]*types.Command, error) {
	var commands []*types.Command
	// Разделение на команды по конвейеру
	parts := strings.Split(input, "|")
	for _, part := range parts {
		part = strings.TrimSpace(part)

		var cmd *types.Command

		// Обработка перенаправлений
		if strings.Contains(part, ">") {
			subParts := strings.SplitN(part, ">", 2)
			cmd = parseCommand(subParts[0])
			cmd.OutputRedirect = strings.TrimSpace(subParts[1])
			if strings.HasPrefix(part, ">>") {
				cmd.Append = true
				cmd.OutputRedirect = strings.TrimSpace(strings.TrimPrefix(subParts[1], ">"))
			}
		} else if strings.Contains(part, "<") {
			subParts := strings.SplitN(part, "<", 2)
			cmd = parseCommand(subParts[0])
			cmd.InputRedirect = strings.TrimSpace(subParts[1])
		} else {
			cmd = parseCommand(part)
		}

		if cmd.Name != "" {
			commands = append(commands, cmd)
		}
	}
	// Связывание команд в конвейер
	for i := 0; i < len(commands)-1; i++ {
		commands[i].PipeTo = commands[i+1]
	}
	return commands, nil
}

func parseCommand(part string) *types.Command {
	args := strings.Fields(part)
	if len(args) == 0 {
		return &types.Command{}
	}
	return &types.Command{
		Name: args[0],
		Args: args[1:],
	}
}
