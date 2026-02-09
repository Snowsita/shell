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
	if strings.Contains(command, string(os.PathSeparator)) {
		info, err := os.Stat(command)
		if err == nil && !info.IsDir() && info.Mode()&0111 != 0 {
			return command
		}
		return ""
	}

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

func parseInput(input string) []string {
	var parts []string
	var currentPart strings.Builder
	inSingleQuotes := false
	inDoubleQuotes := false
	isEscaped := false

	for _, char := range input {

		if isEscaped {
			if inDoubleQuotes && !(char == '"' || char == '\\' || char == '$' || char == '`') {
				currentPart.WriteRune('\\')
			}

			currentPart.WriteRune(char)
			isEscaped = false
			continue
		}

		if char == '\\' && !inSingleQuotes {
			isEscaped = true
			continue
		}

		if char == '\'' && !inDoubleQuotes {
			inSingleQuotes = !inSingleQuotes
			continue
		}

		if char == '"' && !inSingleQuotes {
			inDoubleQuotes = !inDoubleQuotes
			continue
		}

		if char == ' ' && !inSingleQuotes && !inDoubleQuotes {
			if currentPart.Len() > 0 {
				parts = append(parts, currentPart.String())
				currentPart.Reset()
			}
			continue
		}

		currentPart.WriteRune(char)
	}

	if currentPart.Len() > 0 {
		parts = append(parts, currentPart.String())
	}

	return parts
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

		parts := parseInput(input)

		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		info := parseRediretions(parts[1:])

		switch command {
		case "exit":
			os.Exit(0)
		case "echo":
			outW, _ := GetOutputWriter(info.StdoutFile, false, os.Stdout)
			if info.AppendFile != "" {
				outW, _ = GetOutputWriter(info.AppendFile, true, os.Stdout)
			}

			if info.StderrFile != "" {
				errW, err := GetOutputWriter(info.StderrFile, false, os.Stderr)
				if err == nil && errW != os.Stderr {
					errW.Close()
				}
			}

			fmt.Fprintln(outW, strings.Join(info.FinalArgs, " "))
		case "type":
			if len(parts) > 1 {
				target := parts[1]

				if target == "exit" || target == "echo" || target == "type" || target == "pwd" || target == "cd" {
					fmt.Printf("%s is a shell builtin\n", target)
					continue
				}

				fullPath := getExecutablePath(target)
				if fullPath != "" {
					fmt.Printf("%s is %s\n", target, fullPath)
				} else {
					fmt.Printf("%s: not found\n", target)
				}
			}
		case "pwd":
			outW, _ := GetOutputWriter(info.StdoutFile, false, os.Stdout)
			if info.AppendFile != "" {
				outW, _ = GetOutputWriter(info.AppendFile, true, os.Stdout)
			}

			dir, err := os.Getwd()
			if err == nil {
				fmt.Fprintln(outW, dir)
			}

			if outW != os.Stdout {
				outW.Close()
			}
		case "cd":
			if len(parts) < 2 {
				continue
			}

			target := parts[1]
			if target == "~" {
				target = os.Getenv("HOME")
			}

			err := os.Chdir(target)
			if err != nil {
				fmt.Printf("cd: %s: No such file or directory\n", target)
			}

		default:
			fullPath := getExecutablePath(command)

			if fullPath != "" {
				cmd := exec.Command(fullPath, info.FinalArgs...)

				cmd.Args[0] = command

				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				outWriter, _ := GetOutputWriter(info.StdoutFile, false, os.Stdout)
				if info.AppendFile != "" {
					outWriter, _ = GetOutputWriter(info.AppendFile, true, os.Stdout)
				}
				cmd.Stdout = outWriter

				errWriter, _ := GetOutputWriter(info.StderrFile, false, os.Stderr)
				if info.AppendErrFile != "" {
					errWriter, _ = GetOutputWriter(info.AppendErrFile, true, os.Stderr)
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
