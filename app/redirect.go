package main

import "os"

type RedirectInfo struct {
	StdoutFile string
	StderrFile string
	AppendFile string
	AppendErrFile string
	FinalArgs  []string
}

func parseRediretions(args []string) RedirectInfo {
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

func GetOutputWriter(fileName string, isAppend bool, defaultWriter *os.File) (*os.File, error) {
	if fileName == "" {
		return defaultWriter, nil
	}

	flags := os.O_WRONLY | os.O_CREATE
	if isAppend {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	return os.OpenFile(fileName, flags, 0644)
}