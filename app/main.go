package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/shell"
	"os"
	"strings"
)

var _ = fmt.Print

func main() {
	var history []string

	completer := &shell.BuiltinCompleter{
		Builtins: []string{"exit", "echo", "type", "pwd", "cd", "history"},
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "$ ",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
	})
	if err != nil {
		panic(err)
	}

	histFile := os.Getenv("HISTFILE")
	if histFile != "" {
		shell.FileHistory(&history, histFile)
	}

	defer rl.Close()

	for {
		input, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				continue
			}
			break
		}

		input = strings.TrimSpace(input)

		if input != "exit" {
			history = append(history, input)
		}
		
		parts := ParseInput(input)

		if len(parts) == 0 {
			continue
		}

		pipeIndex := -1
		for i, p := range parts {
			if p == "|" {
				pipeIndex = i
				break
			}
		}

		if pipeIndex != -1 {
			runPipeline(&history, parts)
		} else {
			runSingleCommand(&history, parts)
		}
	}
}