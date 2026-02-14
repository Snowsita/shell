package main

import (
	"fmt"
	"github.com/Snowsita/shell/app/shell"
	"os"
	"os/exec"
)

func runBuiltin(history *[]string, name string, info shell.RedirectInfo, out *os.File) {
	switch name {
	case "exit":
		histFile := os.Getenv("HISTFILE")
		if histFile != "" {
			_ = shell.AppendHistory(history, histFile)
		}
		os.Exit(0)

	case "echo":
		shell.HandleEcho(info, out)

	case "type":
		shell.HandleType(info, out, getExecutablePath)

	case "pwd":
		shell.HandlePwd(info, out)

	case "cd":
		if err := shell.HandleCd(info.FinalArgs); err != nil {
			fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", info.FinalArgs[0])
		}

	case "history":
		shell.HandleHistory(history, info, out)
	}
}

func runSingleCommand(history *[]string, parts []string) {
	command := parts[0]
	info := shell.ParseRedirections(parts[1:])

	switch command {
	case "exit", "echo", "type", "pwd", "cd", "history":
		runBuiltin(history, command, info, os.Stdout)
		return
	}

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