package shell

import (
	"os"
)

func HandleCd(parts []string) error {
	if len(parts) == 0 {
		home := os.Getenv("HOME")
		if home == "" {
            home = os.Getenv("USERPROFILE") 
        }
		return os.Chdir(home)
	}

	target := parts[0]
	if target == "~" {
		target = os.Getenv("HOME")
	}

	err := os.Chdir(target)
	if err != nil {
		return err
	}

	return nil
}
