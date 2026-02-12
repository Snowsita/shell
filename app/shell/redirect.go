package shell

import (
	"os"
	"io"
)

type RedirectInfo struct {
	StdoutFile    string
	StderrFile    string
	AppendFile    string
	AppendErrFile string
	FinalArgs     []string
}

func ParseRedirections(args []string) RedirectInfo {
	res := RedirectInfo{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case ">", "1>":
			if i+1 < len(args) {
				res.StdoutFile = args[i+1]
				i++
			}
		case "2>":
			if i+1 < len(args) {
				res.StderrFile = args[i+1]
				i++
			}
		case ">>", "1>>":
			if i+1 < len(args) {
				res.AppendFile = args[i+1]
				i++
			}
		case "2>>":
			if i+1 < len(args) {
				res.AppendErrFile = args[i+1]
				i++
			}
		default:
			res.FinalArgs = append(res.FinalArgs, args[i])
		}
	}
	return res
}

func GetOutputWriter(fileName string, isAppend bool, defaultOut io.Writer) (io.Writer, error) {
	if fileName == "" {
		return defaultOut, nil
	}

	flags := os.O_WRONLY | os.O_CREATE
	if isAppend {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	return os.OpenFile(fileName, flags, 0644)
}

func (info RedirectInfo) GetStdout(defaultOut io.Writer) (io.Writer, error) {
    if info.AppendFile != "" {
        return GetOutputWriter(info.AppendFile, true, defaultOut)
    }
    return GetOutputWriter(info.StdoutFile, false, defaultOut)
}

func (info RedirectInfo) GetStderr(defaultOut io.Writer) (io.Writer, error) {
    if info.AppendErrFile != "" {
        return GetOutputWriter(info.AppendErrFile, true, defaultOut)
    }
    return GetOutputWriter(info.StderrFile, false, defaultOut)
}
