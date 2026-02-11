package shell

import (
	"fmt"
	"strings"
)

type BuiltinCompleter struct {
	Builtins []string
}

func (c *BuiltinCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	var matches [][]rune
	input := string(line[:pos])

	if strings.Contains(input, " ") {
		return nil, 0
	}

	for _, b := range c.Builtins {
		if strings.HasPrefix(b, input) {
			completion := b[len(input):] + " "
			matches = append(matches, []rune(completion))
		}
	}

	if len(matches) == 0 {
		fmt.Print("\x07")
		return nil, 0
	}

	return matches, len(input)
}