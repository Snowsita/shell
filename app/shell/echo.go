package shell

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func HandleEcho(info RedirectInfo, defaultOut io.Writer) error {
	outW, err := info.GetStdout(defaultOut)
	if err != nil {
		return err
	}

	if outW != defaultOut {
		if closer, ok := outW.(io.Closer); ok {
			defer closer.Close()
		}
	}

	errW, err := info.GetStderr(os.Stderr)
	if err != nil {
		return err
	}

	if errW != os.Stderr {
		if closer, ok := errW.(io.Closer); ok {
			defer closer.Close()
		}
	}

	_, err = fmt.Fprintln(outW, strings.Join(info.FinalArgs, " "))
	return err
}
