package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func HandleHistory(history *[]string, info RedirectInfo, defaultOut io.Writer) error {
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

	if len(args) > 0 {
		if args[0] == "-r" {
			if len(args) < 2 {
				return fmt.Errorf("history: argument required")
			}
			err := fileHistory(history, args[1])
			if err != nil {
				return err
			}
			return nil
		}
	}

	startIndex := 0
	if len(args) > 0 {
		n, err := strconv.Atoi(args[0])

		if err == nil {
			if n < len(*history) {
				startIndex = len(*history) - n
			}
		}
	}

	hist := *history
	for i := startIndex; i < len(hist); i++ {
		cmd := hist[i]
		_, err = fmt.Fprintf(outW, "%5d %s\n", i+1, cmd)

		if err != nil {
			return err
		}
	}
	return nil
}

func fileHistory(history *[]string, filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		*history = append(*history, line)
	}

	if err := scanner.Err(); err != nil {
        return err
    }

	return nil
}
