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

func (a *App) AddCommand(cmd *Command) {
	a.Commands = append(a.Commands, cmd)
}
func (a *App) Run() {
	args := os.Args[1:]

	if len(args) > 0 {
		if args[0] == "--help" || args[0] == "-h" {
			a.PrintGlobalHelp()
			return
		}
		if args[0] == "--version" || args[0] == "-v" {
			fmt.Printf("%s version %s\n", a.Name, a.Version)
			return
		}

		for _, cmd := range a.Commands {

			if cmd.Use == args[0] {

				os.Args = os.Args[1:]
				cmd.Execute()
				return
			}
		}
	}

	if a.RootCmd != nil {
		a.RootCmd.Execute()
	} else {
		fmt.Printf("Error: unknown command %q\nRun '%s --help' for usage.\n", args[0], a.Name)
		os.Exit(1)
	}
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

func (a *App) PrintGlobalHelp() {
	fmt.Printf("%s - %s\n\n", a.Name, a.Description)
	fmt.Printf("Usage:\n  %s [command] [flags]\n\n", a.Name)

	if len(a.Commands) > 0 {
		fmt.Println("Available Commands:")
		for _, cmd := range a.Commands {
			fmt.Printf("  %-12s %s\n", cmd.Use, cmd.Short)
		}
		fmt.Println()
	}
	fmt.Printf("Use \"%s [command] --help\" for more information about a command.\n", a.Name)
}
