package main

import (
	"errors"
	"github.com/spf13/cobra"
	"os"

	"github.com/tendermint/tmlibs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"

	coolcmd "github.com/cosmos/cosmos-sdk/examples/ethermint/x/cool/commands"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/commands"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/commands"
	ibccmd "github.com/cosmos/cosmos-sdk/x/ibc/commands"
	stakingcmd "github.com/cosmos/cosmos-sdk/x/staking/commands"

	"github.com/cosmos/cosmos-sdk/examples/ethermint/app"
	"github.com/cosmos/cosmos-sdk/examples/ethermint/types"
)

// ethercliCmd is the entry point for this binary
var (
	ethercliCmd = &cobra.Command{
		Use:   "ethercli",
		Short: "Ethermint light-client",
	}
)

func todoNotImplemented(_ *cobra.Command, _ []string) error {
	return errors.New("TODO: Command not yet implemented")
}

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// get the codec
	cdc := app.MakeCodec()

	// TODO: setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc

	// add standard rpc, and tx commands
	rpc.AddCommands(ethercliCmd)
	ethercliCmd.AddCommand(client.LineBreak)
	tx.AddCommands(ethercliCmd, cdc)
	ethercliCmd.AddCommand(client.LineBreak)

	// add query/post commands (custom to binary)
	ethercliCmd.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("main", cdc, types.GetAccountDecoder(cdc)),
		)...)
	ethercliCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
		)...)
	ethercliCmd.AddCommand(
		client.PostCommands(
			coolcmd.QuizTxCmd(cdc),
		)...)
	ethercliCmd.AddCommand(
		client.PostCommands(
			coolcmd.SetTrendTxCmd(cdc),
		)...)
	ethercliCmd.AddCommand(
		client.PostCommands(
			ibccmd.IBCTransferCmd(cdc),
		)...)
	ethercliCmd.AddCommand(
		client.PostCommands(
			ibccmd.IBCRelayCmd(cdc),
			stakingcmd.BondTxCmd(cdc),
		)...)
	ethercliCmd.AddCommand(
		client.PostCommands(
			stakingcmd.UnbondTxCmd(cdc),
		)...)

	// add proxy, version and key info
	ethercliCmd.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(cdc),
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(ethercliCmd, "EM", os.ExpandEnv("$HOME/.etherli"))
	executor.Execute()
}
