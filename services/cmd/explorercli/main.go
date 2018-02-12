package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/tmlibs/cli"

	"github.com/cosmos/cosmos-sdk/client/commands"
	"github.com/ly0129ly/explorer/services/version"
	"github.com/ly0129ly/explorer/services/modules/sync"
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
	ExplorerCli.PersistentFlags().String(sync.FlagSyncJson, "./sync.json", "json file to save progress")
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
		syncCmd,
		
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(ExplorerCli, "EX", os.ExpandEnv("$HOME/.explorer-cli"))
	executor.Execute()
}
