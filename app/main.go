package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/shell"
	"os"
	"os/exec"
	"strings"
)

var _ = fmt.Print

func main() {
	completer := &shell.BuiltinCompleter{
		Builtins: []string{"exit", "echo", "type", "pwd", "cd"},
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "$ ",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		input, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				continue
			}
		}

		input = strings.TrimSpace(input)

		parts := ParseInput(input)

		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		info := shell.ParseRedirections(parts[1:])

		switch command {
		case "exit":
			os.Exit(0)
		case "echo":
			shell.HandleEcho(info)
		case "type":
			shell.HandleType(info, getExecutablePath)
		case "pwd":
			shell.HandlePwd(info)
		case "cd":
			if err := shell.HandleCd(info.FinalArgs); err != nil {
				fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", info.FinalArgs[0])
			}
		default:
			fullPath := getExecutablePath(command)

			if fullPath != "" {
				cmd := exec.Command(fullPath, info.FinalArgs...)

				cmd.Args[0] = command

				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				outWriter, _ := shell.GetOutputWriter(info.StdoutFile, false, os.Stdout)
				if info.AppendFile != "" {
					outWriter, _ = shell.GetOutputWriter(info.AppendFile, true, os.Stdout)
				}
				cmd.Stdout = outWriter

				errWriter, _ := shell.GetOutputWriter(info.StderrFile, false, os.Stderr)
				if info.AppendErrFile != "" {
					errWriter, _ = shell.GetOutputWriter(info.AppendErrFile, true, os.Stderr)
				}
				cmd.Stderr = errWriter

				cmd.Run()
			} else {
				fmt.Printf("%s: command not found\n", command)
			}
		}
	}
}
