package shell

import (
	"os"
	"errors"
)

func HandleCd(parts []string) error {
	if len(parts) < 2 {
		return errors.New("cd: missing argument")
	}

	target := parts[1]
	if target == "~" {
		target = os.Getenv("HOME")
	}

	err := os.Chdir(target)
	if err != nil {
		return err
	}

	return nil
}
