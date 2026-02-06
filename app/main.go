package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		input = input[:len(input)-1]

		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue
		}

		command := parts[0]

		switch command {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(strings.Join(parts[1:], " "))
		case "type":
			if len(parts) > 1 {
				target := parts[1]

				switch target {
				case "exit", "echo", "type":
					fmt.Printf("%s is a shell builtin\n", target)
				default:
					fmt.Printf("%s not found\n", target)
				}
			}
		default:
			fmt.Printf("%s: command not found\n", command)
		}

	}

}
