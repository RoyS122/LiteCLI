# LiteCLI

`LiteCLI` is an ultra-lightweight and high-performance command-line interface (CLI) engine written in Go. It was designed to replace heavier frameworks like Cobra by eliminating unnecessary memory allocations and massive map structures running in the background.

With `LiteCLI`, your final binaries shrink significantly while maintaining a smooth, typed API and clean routing for subcommands.

---

## Features

* **Zero Unnecessary Allocations**: No heavy background data structures, just fast static typing.
* **Application Routing (`App`)**: Native handling of subcommands (e.g., `app process [flags]`) with complete isolation of positional arguments.
* **Direct Binding via Pointers**: No need to retrieve flag values using string keys. You pass the pointers of your variables (`&quality`, `&output`) directly.
* **Memory Efficiency**: Uses `rune` for short flags to prevent useless sub-string allocations.
* **Automatic Global Help**: Automatically generates help and version menus from your registered commands.

---

## Supported Flag Types

`LiteCLI` handles direct pointer binding for the following standard Go types:

| Method | Go Type | CLI Syntax Example |
| --- | --- | --- |
| `StringVarP` | `string` | `--output="./dist"` or `-o "./dist"` |
| `IntVarP` | `int` | `--quality=80` or `-q 80` |
| `Uint8VarP` | `uint8` | `--workers=4` or `-w 4` |
| `Uint64VarP` | `uint64` | `--max-size=5000000000` |
| `Float64VarP` | `float64` | `--ratio=1.5` or `-r 1.5` |
| `BoolVarP` | `bool` | `--verbose` or `-v` *(boolean switch, no value required)* |

---

## Installation

Add `LiteCLI` to your Go project:

```bash
go get github.com/RoyS122/LiteCLI@v1.0.1

```

---

## Architecture & Usage

The architecture cleanly separates the application (`App`), which handles global routing, from the commands (`Command`), which hold local flags and execution logic.

### 1. Declaring a Command (e.g., `cmd/process.go`)

```go
package cmd

import (
	"fmt"
	"github.com/RoyS122/LiteCLI"
)

var (
	outputPath string
	quality    int
)

var ProcessCmd = &litecli.Command{
	Name:  "process",
	Short: "Runs the processing or compression",
	Long:  "A longer description of the process command...",
	Run: func(cmd *litecli.Command, args []string) {
		fmt.Printf("Target directory: %v\n", args)
		fmt.Printf("Flags -> Output: %s | Quality: %d%%\n", outputPath, quality)
	},
}

func init() {
	// Direct flag binding via pointers (LiteCLI's fluid API)
	ProcessCmd.StringVarP(&outputPath, "output", 'o', "./dist")
	ProcessCmd.IntVarP(&quality, "quality", 'q', 80)
}

```

### 2. Application Entry Point (e.g., `cmd/root.go`)

```go
package cmd

import (
	"github.com/RoyS122/LiteCLI"
)

var RootCmd = &litecli.Command{
	Short: "Main tool",
	Run: func(cmd *litecli.Command, args []string) {
		println("Welcome! Specify a command, for example 'process'.")
	},
}

func Execute() {
	// Centralizing the entry point with the App object
	app := litecli.NewApp("my_tool", "An ultra-fast CLI", "1.0.1", RootCmd)

	// Registering subcommands
	app.AddCommand(ProcessCmd)

	// Launching parsing and routing
	app.Run()
}

```

---

## License

This project is licensed under the **MIT License**. You are free to use, modify, and distribute it.
