package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/tmlibs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/ibc"
)

var (
	ibcCmd = &cobra.Command{
		Use:   "ibc",
		Short: "IBC command-line interface",
	}
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := ibc.MakeCodec()

	ibcCmd.AddCommand(
		sendCmd,
		relayCmd,
	)

	executor := cli.PrepareMainCmd(ibcCmd, "IBC", os.ExpandEnv("$HOME/.ibccli"))
	executor.Execute()
}
