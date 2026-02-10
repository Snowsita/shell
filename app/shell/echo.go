package shell

import (
	"fmt"
	"os"
	"strings"
)

func HandleEcho(info RedirectInfo) error {
	outW, err := info.GetStdout(os.Stdout)
	if err != nil {
		return err
	}

	if outW != os.Stdout {
		defer outW.Close()
	}

	errW, err := info.GetStderr(os.Stderr)
	if err != nil {
		return err
	}

	if errW != os.Stderr {
		defer errW.Close()
	}

	_, err = fmt.Fprintln(outW, strings.Join(info.FinalArgs, " "))
	return err
}
