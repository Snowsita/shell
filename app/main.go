package main

import (
	"bufio"
	"fmt"
	"github.com/codecrafters-io/shell-starter-go/app/shell"
	"os"
	"os/exec"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
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

				if f, ok := cmd.Stdout.(*os.File); ok && f != os.Stdout {
					f.Close()
				}
				if f, ok := cmd.Stderr.(*os.File); ok && f != os.Stderr {
					f.Close()
				}

			} else {
				fmt.Printf("%s: command not found\n", command)
			}
		}
	}
}
