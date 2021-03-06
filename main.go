package main

import (
	"fmt"
	"os"

	"github.com/jrasell/levant/buildtime"
	"github.com/mitchellh/cli"
)

var (
	// variables populated by govvv(1)
	Version    = "dev"
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
)

func main() {
	exportBuildtimeConsts()
	os.Exit(Run(os.Args[1:]))
}

// Run sets up the commands and triggers RunCustom which inacts the correct
// run of Levant.
func Run(args []string) int {
	return RunCustom(args, Commands(nil))
}

// RunCustom is the main function to trigger a run of Levant.
func RunCustom(args []string, commands map[string]cli.CommandFactory) int {
	// Get the command line args. We shortcut "--version" and "-v" to
	// just show the version.
	for _, arg := range args {
		if arg == "-v" || arg == "-version" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	// Build the commands to include in the help now.
	commandsInclude := make([]string, 0, len(commands))
	for k := range commands {
		switch k {
		default:
			commandsInclude = append(commandsInclude, k)
		}
	}

	cli := &cli.CLI{
		Args:     args,
		Commands: commands,
		HelpFunc: cli.FilteredHelpFunc(commandsInclude, cli.BasicHelpFunc("levant")),
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}

func exportBuildtimeConsts() {
	buildtime.GitCommit = GitCommit
	buildtime.GitBranch = GitBranch
	buildtime.GitState = GitState
	buildtime.GitSummary = GitSummary
	buildtime.BuildDate = BuildDate
	buildtime.Version = Version
}
