package executor

import (
	"fmt"
	"os"
	"os/exec"
	"wb_tech/l2_15/pkg/types"
)

func Execute(cmd *types.Command) error {

	switch cmd.Name {
	case "cd":
		return cd(cmd.Args)

	case "pwd":
		return pwd(cmd.Args)
	default:
		return Internalexec(cmd.Name, cmd.Args)
	}
}
func cd(args []string) error {
	fmt.Println("cd")
	return nil
}

func pwd(args []string) error {
	return nil
}
func Internalexec(name string, args []string) error {
	fmt.Println("Internalexec")
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
