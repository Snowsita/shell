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
    var found []string
    input := string(line[:pos])

    if strings.Contains(input, " ") {
        return nil, 0
    }

    // 1. Gather all (Builtins + Path)
    for _, b := range c.Builtins {
        if strings.HasPrefix(b, input) {
            found = append(found, b)
        }
    }
    external := FindPathMatches(input)
    found = append(found, external...)
    sort.Strings(found)

    if len(found) == 0 {
        fmt.Print("\x07")
        return nil, 0
    }

    // 2. Handle Multiple Matches (The Manual Way)
    if len(found) > 1 {
        fmt.Print("\x07") // Ring the bell as requested
        
        // This is the magic part:
        // We print a newline, the matches joined by THREE spaces, 
        // a newline, and then we RESTORE the prompt line.
        fmt.Printf("\n%s\n$ %s", strings.Join(found, "  "), input)
        
        // Return nil so the library doesn't try to print its own 1-space version
        return nil, 0
    }

    // 3. Handle Single Match
    match := found[0]
    completion := match[len(input):] + " "
    return [][]rune{[]rune(completion)}, len(input)
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
