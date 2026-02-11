package shell

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type BuiltinCompleter struct {
	Builtins []string
	TabCount int
}

func (c *BuiltinCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	var allMatches []string
	input := string(line[:pos])

	if strings.Contains(input, " ") {
		c.TabCount = 0
		return nil, 0
	}

	for _, b := range c.Builtins {
		if strings.HasPrefix(b, input) {
			allMatches = append(allMatches, b)
		}
	}

	externalMatches := FindPathMatches(input)
	allMatches = append(allMatches, externalMatches...)

	sort.Strings(allMatches)

	if len(allMatches) == 0 {
		fmt.Print("\x07")
		c.TabCount = 0
		return nil, 0
	}

	if len(allMatches) == 1 {
		c.TabCount = 0
		completion := allMatches[0] + " "
		return [][]rune{[]rune(completion)}, 0
	}

	if len(allMatches) > 1 {
		c.TabCount++

		if c.TabCount == 1 {
			fmt.Print("\x07")
			return nil, 0
		}

		fmt.Printf("\n%s\n$ %s", strings.Join(allMatches, "  "), input)
		c.TabCount = 0

		return nil, 0
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
