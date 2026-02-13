package shell

import (
	"fmt"
	"io"
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

	for i, cmd := range history {
		_, err = fmt.Fprintf(outW, "%5d %s\n", i+1, cmd)

		if err != nil {
			return err
		}
	}
	return nil
}
