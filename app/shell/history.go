package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

var historyWrittenIndex = 0

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
		switch args[0] {
		case "-r":
			if len(args) < 2 {
				return fmt.Errorf("history: argument required")
			}
			return fileHistory(history, args[1])
		case "-w":
			if len(args) < 2 {
				return fmt.Errorf("history: argument required")
			}
			return writeHistory(history, args[1])
		case "-a":
			if len(args) < 2 {
				return fmt.Errorf("history: argument required")
			}
			return appendHistory(history, args[1])
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

func writeHistory(history *[]string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	hist := *history

	for _, cmd := range hist {
		_, err := fmt.Fprintln(file, cmd)

		if err != nil {
			return err
		}
	}

	historyWrittenIndex = len(*history)

	return nil
}

func appendHistory(history *[]string, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	hist := *history

	if historyWrittenIndex < len(hist) {
		for _, cmd := range hist[historyWrittenIndex:] {
			if _, err := fmt.Fprintln(file, cmd); err != nil {
				return err
			}
		}
        
        historyWrittenIndex = len(hist)
	}

	return nil
}
