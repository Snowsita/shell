package shell

import (
	"fmt"
	"os"
)

func HandlePwd(info RedirectInfo) error {
	outW, _ := GetOutputWriter(info.StdoutFile, false, os.Stdout)
	if info.AppendFile != "" {
		outW, _ = GetOutputWriter(info.AppendFile, true, os.Stdout)
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Fprintln(outW, dir)

	if outW != os.Stdout {
		outW.Close()
	}

	return nil
}
