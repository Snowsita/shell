package shell

import (
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

	return matches, len(input)
}