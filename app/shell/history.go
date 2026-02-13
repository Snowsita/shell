package shell

import (
	"fmt"
	"io"
	"strconv"
)

func HandleHistory(history []string, info RedirectInfo, defaultOut io.Writer) error {
	outW, err := info.GetStdout(defaultOut)
	if err != nil {
		return err
	}

	if outW != defaultOut {
		if closer, ok := outW.(io.Closer); ok {
			defer closer.Close()
		}
	}

	args := info.FinalArgs
	startIndex := 0

	if len(args) > 0 {
		n, err := strconv.Atoi(args[0])

		if err == nil {
			if n < len(history) {
				startIndex = len(history) - n
			}
		}
	}

	for i := startIndex; i < len(history); i++ {
		cmd := history[i]
		_, err = fmt.Fprintf(outW, "%5d %s\n", i+1, cmd)

		if err != nil {
			return err
		}
	}
	return nil
}
