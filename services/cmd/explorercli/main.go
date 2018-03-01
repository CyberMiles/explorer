package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/tmlibs/cli"

	"github.com/cosmos/cosmos-sdk/client/commands"
	"github.com/cybermiles/explorer/services/version"
)

// entry point for this binary
var (
	ExplorerCli = &cobra.Command{
		Use:   "explorercli",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func prepareMainCmd() {
	commands.AddBasicFlags(ExplorerCli)
}

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// add commands
	prepareMainCmd()
	prepareRestServerCommands()

	ExplorerCli.AddCommand(
		commands.InitCmd,
		restServerCmd,
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(ExplorerCli, "EX", os.ExpandEnv("$HOME/.explorer-cli"))
	executor.Execute()
}
