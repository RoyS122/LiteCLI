package litecli

import (
	"fmt"
	"os"
)

func NewApp(name, description, version string, rootCmd *Command) *App {
	return &App{
		Name:        name,
		Description: description,
		Version:     version,
		RootCmd:     rootCmd,
	}
}

func (a *App) Run() {
	if a.RootCmd == nil {
		fmt.Println("Error: no root command defined for this app")
		os.Exit(1)
	}

	// Intercept de base pour le help universel ou la version
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "--help" || arg == "-h" {
			a.RootCmd.PrintHelp(a.Name)
			return
		}
		if arg == "--version" || arg == "-v" {
			fmt.Printf("%s version %s\n", a.Name, a.Version)
			return
		}
	}

	a.RootCmd.Execute()
}

func (c *Command) PrintHelp(appName string) {
	fmt.Printf("Usage:\n  %s [flags] [arguments]\n\n", appName)
	if c.Long != "" {
		fmt.Printf("%s\n\n", c.Long)
	} else if c.Short != "" {
		fmt.Printf("%s\n\n", c.Short)
	}

	if len(c.Flags) > 0 {
		fmt.Println("Flags:")
		for _, f := range c.Flags {
			var shortStr string
			if f.Short != 0 {
				shortStr = fmt.Sprintf("-%c, ", f.Short)
			} else {
				shortStr = "    "
			}
			fmt.Printf("  %s--%-12s (default: %s)\n", shortStr, f.Name, f.DefaultValue)
		}
		fmt.Println()
	}
}
