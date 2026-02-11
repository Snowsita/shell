package shell

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type BuiltinCompleter struct {
	Builtins []string
	LastInput string
	TabCount int
}

func (c *BuiltinCompleter) Do(line []rune, pos int) ([][]rune, int) {
	input := string(line[:pos])

	if strings.Contains(input, " ") {
		c.TabCount = 0
		return nil, 0
	}

	var matches []string

	for _, b := range c.Builtins {
		if strings.HasPrefix(b, input) {
			matches = append(matches, b)
		}
	}

	externalMatches := FindPathMatches(input)
	matches = append(matches, externalMatches...)

	sort.Strings(matches)

	if len(matches) == 0 {
		c.TabCount = 0
		return nil, 0
	}

	// Reset tab counter if input changed
	if input != c.LastInput {
		c.TabCount = 0
	}
	c.LastInput = input
	c.TabCount++

	// First TAB
	if c.TabCount == 1 && len(matches) > 1 {
		fmt.Print("\x07")
		return nil, 0
	}

	// Second TAB
	if c.TabCount >= 2 && len(matches) > 1 {
		fmt.Print("\n")
		fmt.Print(strings.Join(matches, "  "))
		fmt.Print("\n$ ")
		fmt.Print(input)

		c.TabCount = 0
		return nil, 0
	}

	// Single match
	if len(matches) == 1 {
		suffix := matches[0][len(input):] + " "
		c.TabCount = 0
		return [][]rune{[]rune(suffix)}, len(input)
	}

	return nil, 0
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
