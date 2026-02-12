package main

import (
	"strings"
	"os"
	"path/filepath"
)

type ParserInput struct {
	Parts []string
	CurrentPart strings.Builder
	InSingleQuotes bool
	InDoubleQuotes bool
	IsEscaped bool
}

func ParseInput(input string) []string {
	res := ParserInput{}

	for _, char := range input {

		if res.IsEscaped {
			if res.InDoubleQuotes && !(char == '"' || char == '\\' || char == '$' || char == '`') {
				res.CurrentPart.WriteRune('\\')
			}

			res.CurrentPart.WriteRune(char)
			res.IsEscaped = false
			continue
		}

		if char == '\\' && !res.InSingleQuotes {
			res.IsEscaped = true
			continue
		}

		if char == '\'' && !res.InDoubleQuotes {
			res.InSingleQuotes = !res.InSingleQuotes
			continue
		}

		if char == '"' && !res.InSingleQuotes {
			res.InDoubleQuotes = !res.InDoubleQuotes
			continue
		}

		if char == ' ' && !res.InSingleQuotes && !res.InDoubleQuotes {
			if res.CurrentPart.Len() > 0 {
				res.Parts = append(res.Parts, res.CurrentPart.String())
				res.CurrentPart.Reset()
			}
			continue
		}

		res.CurrentPart.WriteRune(char)
	}

	if res.CurrentPart.Len() > 0 {
		res.Parts = append(res.Parts, res.CurrentPart.String())
	}

	return res.Parts
}

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

func isBuiltin(cmd string) bool {
	switch cmd {
	case "echo", "type", "pwd", "exit", "cd":
		return true
	}
	return false
}