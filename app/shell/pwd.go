package shell

import (
	"fmt"
	"os"
	"io"
)

func HandlePwd(info RedirectInfo, defaultOut io.Writer) error {
	outW, _ := GetOutputWriter(info.StdoutFile, false, defaultOut)
	if info.AppendFile != "" {
		outW, _ = GetOutputWriter(info.AppendFile, true, defaultOut)
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Fprintln(outW, dir)

	if outW != defaultOut {
		if closer, ok := outW.(io.Closer); ok {
			defer closer.Close()
		}
	}

	return nil
}
