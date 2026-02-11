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
    // GLOBAL DEDUPLICATION MAP (Crucial for echo vs /bin/echo)
    seen := make(map[string]bool) 
    input := string(line[:pos])

    if strings.Contains(input, " ") {
        c.TabCount = 0
        return nil, 0
    }

    // 1. Gather Matches (with deduplication)
    // Add Builtins
    for _, b := range c.Builtins {
        if strings.HasPrefix(b, input) {
            allMatches = append(allMatches, b)
            seen[b] = true
        }
    }
    // Add External
    externalMatches := FindPathMatches(input)
    for _, ext := range externalMatches {
        if !seen[ext] {
            allMatches = append(allMatches, ext)
            seen[ext] = true
        }
    }

    // 2. Sort
    sort.Strings(allMatches)

    // 3. Handle No Matches
    if len(allMatches) == 0 {
        fmt.Print("\x07")
        c.TabCount = 0
        return nil, 0
    }

    // 4. Handle Single Match (The "Append Only" Fix)
    if len(allMatches) == 1 {
        c.TabCount = 0
        match := allMatches[0]
        
        // STRATEGY: Calculate the Suffix
        // Input: "custom"
        // Match: "custom_exe_2756"
        // Suffix: "_exe_2756 "
        suffix := match[len(input):] + " " 
        
        // Return the SUFFIX with length 0.
        // Length 0 tells readline: "Do not backspace. Just print these characters."
        // This is safer than replacing the whole word.
        return [][]rune{[]rune(suffix)}, 0 // <--- LENGTH MUST BE 0
    }

    // 5. Handle Multiple Matches (The Grid Fix)
    if len(allMatches) > 1 {
        c.TabCount++

        if c.TabCount == 1 {
            fmt.Print("\x07")
            return nil, 0
        }

        // Manual Print with Double Spaces
        formattedList := strings.Join(allMatches, "  ")
        fmt.Printf("\n%s\n$ %s", formattedList, input)
        
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
