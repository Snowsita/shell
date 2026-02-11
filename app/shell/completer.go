package shell

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type BuiltinCompleter struct {
	Builtins []string
}

func (c *BuiltinCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	var matches []string
	input := string(line[:pos])

	if strings.Contains(input, " ") {
		return nil, 0
	}

	for _, b := range c.Builtins {
		if strings.HasPrefix(b, input) {
			matches = append(matches, b)
		}
	}

	externalMatches := FindPathMatches(input)
	matches = append(matches, externalMatches...)

	sort.Strings(matches)

	if len(matches) == 0 {
		fmt.Print("\x07")
		return nil, 0
	}

	if len(matches) > 1 {
		fmt.Print("\x07")
	}

	var finalMatches [][]rune
	for _, match := range matches {
		var completion string
		if len(matches) == 1 {
			completion = match[len(input):] + " "
		} else {
			completion = match[len(input):]
		}
		finalMatches = append(finalMatches, []rune(completion))
	}

	return finalMatches, len(input)
}

func FindPathMatches(prefix string) []string {
	var matches []string
	seen := make(map[string]bool)
	pathEnv := os.Getenv("PATH")
	dirs := strings.SplitSeq(pathEnv, string(os.PathListSeparator))

	for dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			name := entry.Name()

			if strings.HasPrefix(name, prefix) && !seen[name] {
				if isExecutable(entry) {
					matches = append(matches, name)
					seen[name] = true
				}
			}
		}
	}

	return matches
}

func isExecutable(entry os.DirEntry) bool {
	if entry.IsDir() {
		return false
	}

	info, err := entry.Info()
	if err != nil {
		return false
	}

	return info.Mode().Perm()&0111 != 0
}
