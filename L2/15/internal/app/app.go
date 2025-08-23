package app

import (
	"fmt"
	"strings"
	"wb_tech/l2_15/internal/executor"
	"wb_tech/l2_15/internal/parser"
	"wb_tech/l2_15/internal/reader"
)

// Run starts the main loop of the shell.
func Run() {
	for {
		fmt.Print("bsh>: ")
		raw := reader.ReadInput()
		raw = strings.TrimSpace(raw) // убираем \n
		if raw == "" {
			continue
		}
		commands, err := parser.Parse(raw)
		if err != nil {
			fmt.Println("Error parsing command:", err)
			continue
		}
		if len(commands) > 0 {
			executor.Execute(commands)
		}
	}
}
