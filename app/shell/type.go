package shell

import (
	"slices"
	"fmt"
	"io"
)

func HandleType(info RedirectInfo, defaultOut io.Writer, getExecutablePath func(string) string) {
	if len(info.FinalArgs) > 0 {
		target := info.FinalArgs[0]

		if isBuiltin(target) {
			writeOutput(fmt.Sprintf("%s is a shell builtin\n", target))
			return
		}

		fullPath := getExecutablePath(target)
		if fullPath != "" {
			writeOutput(fmt.Sprintf("%s is %s\n", target, fullPath))
		} else {
			writeOutput(fmt.Sprintf("%s: not found\n", target))
		}
	}
}

func isBuiltin(target string) bool {
	builtins := []string{"exit", "echo", "type", "pwd", "cd"}
	return slices.Contains(builtins, target)
}

func writeOutput(message string) {
	fmt.Print(message)
}
