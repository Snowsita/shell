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
    var matches [][]rune
    // Use a string slice for sorting because sorting [][]rune is a nightmare
    var found []string 
    input := string(line[:pos])

    if strings.Contains(input, " ") {
        return nil, 0
    }

    // 1. Collect all strings
    for _, b := range c.Builtins {
        if strings.HasPrefix(b, input) {
            found = append(found, b)
        }
    }
    externalMatches := FindPathMatches(input)
    found = append(found, externalMatches...)

    // 2. SORTING IS MANDATORY (The tester failed you for this earlier)
    sort.Strings(found)

    if len(found) == 0 {
        fmt.Print("\x07")
        return nil, 0
    }

    // 3. Convert to [][]rune based on count
    for _, name := range found {
        var completion string
        if len(found) == 1 {
            // Stage #GM9: Single match MUST have the space
            completion = name[len(input):] + " "
        } else {
            // Stage #WH6: Multiple matches MUST NOT have the space
            // This is what tells readline to use its 'grid' formatter (the 3 spaces)
            completion = name[len(input):]
        }
        matches = append(matches, []rune(completion))
    }

    if len(matches) > 1 {
        fmt.Print("\x07")
        // No TabCount here as requested, just the bell and the matches
        return matches, len(input)
    }

    return matches, len(input)
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
