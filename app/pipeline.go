package main

import (
	"fmt"
	"github.com/codecrafters-io/shell-starter-go/app/shell"
	"os"
	"os/exec"
)

func runPipeline(parts []string) {

	commands := parseCommands(parts)

	var previousPipeReader *os.File = nil

	for i, cmdArgs := range commands {
		isLast := i == len(commands)-1

		info := shell.ParseRedirections(cmdArgs[1:])
		cmdName := cmdArgs[0]

		var nextPipeReader *os.File
		var currentPipeWriter *os.File
		var err error

		if !isLast {
			nextPipeReader, currentPipeWriter, err = os.Pipe()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error creating pipe:", err)
				return
			}
		}

		var cmdStdin *os.File = previousPipeReader
		var cmdStdout *os.File = currentPipeWriter

		if i == 0 && cmdStdin == nil {
			cmdStdin = os.Stdin
		}

		if isLast {
			cmdStdout = os.Stdout
		}

		if isBuiltin(cmdName) {
			if !isLast {
				go func(in *os.File, out *os.File, name string, args shell.RedirectInfo) {
					runBuiltin(name, args, out)
					if out != nil {
						out.Close()
					}
				}(cmdStdin, cmdStdout, cmdName, info)
			} else {
				runBuiltin(cmdName, info, cmdStdout)
			}
		} else {
			cmd := exec.Command(cmdName, info.FinalArgs...)
			cmd.Stdin = cmdStdin
			cmd.Stdout = cmdStdout
			cmd.Stderr = os.Stderr

			if err := cmd.Start(); err != nil {
				fmt.Fprintf(os.Stderr, "Error starting %s: %s\n", cmdName, err)
			}

			if isLast {
				cmd.Wait()
			}
		}

		if currentPipeWriter != nil && !isBuiltin(cmdName) {
			currentPipeWriter.Close()
		}

		if previousPipeReader != nil {
			previousPipeReader.Close()
		}

		previousPipeReader = nextPipeReader
	}
}

func parseCommands(parts []string) [][]string {
	var commands [][]string
	var currentCmd []string

	for _, p := range parts {
		if p == "|" {
			if len(currentCmd) > 0 {
				commands = append(commands, currentCmd)
				currentCmd = []string{}
			}
		} else {
			currentCmd = append(currentCmd, p)
		}
	}

	if len(currentCmd) > 0 {
		commands = append(commands, currentCmd)
	}

	return commands
}
