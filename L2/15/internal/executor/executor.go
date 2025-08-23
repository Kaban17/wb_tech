package executor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"wb_tech/l2_15/pkg/types"
)

// Execute executes a slice of commands, handling pipelines.
func Execute(commands []*types.Command) error {
	if len(commands) == 0 {
		return nil
	}
	if len(commands) == 1 {
		return executeSingleCommand(commands[0])
	}
	return executePipeline(commands)
}

func executeSingleCommand(cmd *types.Command) error {
	var out io.Writer = os.Stdout
	if cmd.OutputRedirect != "" {
		flag := os.O_WRONLY | os.O_CREATE
		if cmd.Append {
			flag |= os.O_APPEND
		} else {
			flag |= os.O_TRUNC
		}
		file, err := os.OpenFile(cmd.OutputRedirect, flag, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		out = file
	}

	switch cmd.Name {
	case "cd":
		return cd(cmd.Args)
	case "pwd":
		return pwd()
	case "echo":
		return echo(cmd.Args, out)
	case "kill":
		return kill(cmd.Args)
	case "ps":
		return ps()
	case "exit":
		os.Exit(0)
	default:
		return externalExec(cmd)
	}
	return nil
}

func cd(args []string) error {
	if len(args) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		return os.Chdir(home)
	}
	return os.Chdir(args[0])
}

func pwd() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(dir)
	return nil
}

func echo(args []string, out io.Writer) error {
	fmt.Fprintln(out, strings.Join(args, " "))
	return nil
}

func kill(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("kill: missing pid")
	}
	pid := 0
	_, err := fmt.Sscan(args[0], &pid)
	if err != nil {
		return fmt.Errorf("kill: invalid pid: %s", args[0])
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Signal(syscall.SIGKILL)
}

func ps() error {
	cmd := exec.Command("ps", "aux")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func externalExec(cmd *types.Command) error {
	execCmd := exec.Command(cmd.Name, cmd.Args...)

	if cmd.InputRedirect != "" {
		file, err := os.Open(cmd.InputRedirect)
		if err != nil {
			return err
		}
		defer file.Close()
		execCmd.Stdin = file
	} else {
		execCmd.Stdin = os.Stdin
	}

	if cmd.OutputRedirect != "" {
		flag := os.O_WRONLY | os.O_CREATE
		if cmd.Append {
			flag |= os.O_APPEND
		} else {
			flag |= os.O_TRUNC
		}
		file, err := os.OpenFile(cmd.OutputRedirect, flag, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		execCmd.Stdout = file
	} else {
		execCmd.Stdout = os.Stdout
	}

	execCmd.Stderr = os.Stderr
	return execCmd.Run()
}

func executePipeline(commands []*types.Command) error {
	var err error
	var cmds []*exec.Cmd
	for _, cmd := range commands {
		cmds = append(cmds, exec.Command(cmd.Name, cmd.Args...))
	}
	for i := 0; i < len(cmds)-1; i++ {
		cmds[i+1].Stdin, err = cmds[i].StdoutPipe()
		if err != nil {
			return err
		}
	}
	cmds[len(cmds)-1].Stdout = os.Stdout
	cmds[0].Stderr = os.Stderr
	for _, cmd := range cmds {
		err = cmd.Start()
		if err != nil {
			return err
		}
	}
	for _, cmd := range cmds {
		err = cmd.Wait()
		if err != nil {
			return err
		}
	}
	return nil
}
