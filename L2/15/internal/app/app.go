package app

import (
	"fmt"
	"strings"
	"wb_tech/l2_15/internal/executor"
	"wb_tech/l2_15/internal/reader"
	"wb_tech/l2_15/pkg/types"
)

func Run() {
	for {
		fmt.Print("bsh>: ")
		raw := reader.ReadInput()
		raw = strings.TrimSpace(raw) // убираем \n
		cmd := types.ToCommand(raw)
		executor.Execute(&cmd)
	}
}
