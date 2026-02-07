package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getExecutablePath(command string) string {
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	for _, dir := range paths {
		fullPath := filepath.Join(dir, command)
		info, err := os.Stat(fullPath)
		if err == nil && !info.IsDir() && info.Mode()&0111 != 0 {
			return fullPath
		}
	}

	return ""
}

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
		args := parts[1:]

		switch command {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(strings.Join(parts[1:], " "))
		case "type":
			if len(parts) > 1 {
				target := parts[1]

				if target == "exit" || target == "echo" || target == "type" || target == "pwd" {
					fmt.Printf("%s is a shell builtin\n", target)
					continue
				}

				pathEnv := os.Getenv("PATH")
				paths := strings.Split(pathEnv, string(os.PathListSeparator))
				found := false

				for _, dir := range paths {
					fullPath := filepath.Join(dir, target)

					info, err := os.Stat(fullPath)
					if err == nil {
						if !info.IsDir() && info.Mode()&0111 != 0 {
							fmt.Printf("%s is %s\n", target, fullPath)
							found = true
							break
						}
					}
				}

				if !found {
					fmt.Printf("%s: not found\n", target)
				}
			}
		case "pwd":
			dir, err := os.Getwd()
			if err == nil {
				fmt.Println(dir)
			} else {
				fmt.Printf("%s: not found", err)
			}
		default:
			fullPath := getExecutablePath(command)

			if fullPath != "" {
				cmd := exec.Command(fullPath, args...)

				cmd.Args[0] = command

				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				err := cmd.Run()
				if err != nil {
					fmt.Printf("%s: error executing command\n", command)
				}
			} else {
				fmt.Printf("%s: command not found\n", command)
			}
		}
	}
}
