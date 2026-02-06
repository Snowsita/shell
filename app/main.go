package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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

		input = input[:len(input)-1]

		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue
		}

		command := parts[0]

		switch command {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(strings.Join(parts[1:], " "))
		case "type":
			if len(parts) > 1 {
				target := parts[1]

				if target == "exit" || target == "echo" || target == "type" {
					fmt.Printf("%s is a shell builtin\n", target)
					continue
				}

				pathEnv := os.Getenv("PATH")
				paths := strings.Split(pathEnv, ":")
				found := false

				for _, dir := range paths {
					fullPath := filepath.Join(dir, target)

					if _, err := os.Stat(fullPath); err == nil {
						fmt.Printf("%s is %s\n", target, fullPath)
						found = true
						break
					}
				}

				if !found {
					fmt.Printf("%s: not found\n", target)
				}
			}
		default:
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
