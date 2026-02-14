# GoShell

A lightweight, performant, and POSIX-compliant shell implementation written entirely in Go.

This project was built from scratch as a deep dive into systems programming, process management, and the Go standard library. It features a custom REPL, command pipelining, I/O redirection, and persistent history management.

## Features

### Core Capabilities
* **REPL (Read-Eval-Print Loop):** fast and responsive interactive prompt using `readline`.
* **External Command Execution:** Executes any binary found in your `$PATH`.
* **Signal Handling:** Gracefully handles `Ctrl+C` (interrupt) and `Ctrl+D` (EOF/Exit).

### Advanced Shell Features
* **Pipelining (`|`):** Chain commands together, passing the output of one as the input to the next.
* **I/O Redirection:**
    * Standard Output: `>` (overwrite), `>>` (append).
    * Standard Error: `2>`, `2>>`.
* **Quoting Support:** Handles single (`'`) and double (`"`) quotes, allowing for preserved whitespace and escaped characters.

### Built-in Commands
The shell includes several internal commands optimized for performance:
* `cd`: Change directory (supports absolute and relative paths).
* `pwd`: Print current working directory.
* `echo`: Print arguments to standard output.
* `type`: Reveal information about commands (builtin vs. external path).
* `exit`: Gracefully terminate the shell session (saves history automatically).
* `history`: View and manage command history.
    * `history`: View session history.
    * `history -w <file>`: Write current history to file.
    * `history -r <file>`: Read history from file.
    * `history -a <file>`: Append new session commands to file.

### Persistence
* **Automatic History Loading:** Loads command history from the file specified in the `HISTFILE` environment variable on startup.
* **Graceful Shutdown:** Automatically appends new commands to `HISTFILE` upon exit.

## Architecture

The project follows a modular "Traffic Controller" architecture:

* **`main.go`**: Acts as the entry point and router. It handles the input loop, parses the raw command string, detects pipes/redirections, and routes execution to either the builtin handler or the OS process spawner.
* **`shell/` package**: Contains the core business logic.
    * **`builtins.go`**: Implementation of internal commands.
    * **`history.go`**: Stateless and session-aware history management.
    * **`parser.go`**: specialized parsing logic for quotes and arguments.

This structure ensures that the execution flow is decoupled from the specific implementation of commands, allowing for easy addition of new features.

## Installation & Usage

### Prerequisites
* Go 1.22 or higher

### Build
```bash
git clone https://github.com/Snowsita/shell.git
cd shell
go build ./app
```
