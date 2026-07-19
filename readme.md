# LiteCLI

A lightweight, zero-dependency, high-performance CLI argument parsing library for Go. Designed as a minimal and blazing fast alternative to heavy frameworks.

## Features

* **Zero Third-Party Dependencies:** Pure Go standard library.
* **Ultra Fast:** Compiles instantly and executes in milliseconds with minimal allocations.
* **Flexible Syntax:** Supports long flags (--flag), short flags (-f), both space-separated and values assigned via =.
* **Type-Safe Binding:** Bind inputs directly to string, int, or uint8 pointers.
* **Automatic Help Generation:** Built-in usage and flag documentation.

## Installation

go get [github.com/RoyS122/LiteCLI](https://www.google.com/search?q=https%3A%2F%2Fgithub.com%2FRoyS122%2FLiteCLI)

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/RoyS122/LiteCLI"
)

func main() {
	var (
		quality uint8
		output  string
	)

	// 1. Define the root command
	rootCmd := &litecli.Command{
		Short: "GoStamp image processing utility",
		Run: func(cmd *litecli.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Error: missing target file")
				return
			}
			fmt.Printf("Processing %s (Quality: %d%%, Target: %s)\n", args[0], quality, output)
		},
	}

	// 2. Bind flags to variables (with automatic default values)
	rootCmd.Uint8VarP(&quality, "quality", 'q', 80)
	rootCmd.StringVarP(&output, "output", 'o', "./dist")

	// 3. Package and run the App
	app := litecli.NewApp("gostamp", "A fast metadata stripping tool", "1.0.0", rootCmd)
	app.Run()
}

```

## Running Tests

go test -v ./...

## License

MIT License - see the LICENSE file for details.
