package commands

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	wire "github.com/tendermint/go-wire"

	"github.com/cosmos/cosmos-sdk/x/ibc"
)

const (
	flagTo     = "to"
	flagAmount = "amount"
)

func IBCTransferCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := commander{cdc}

	cmd := &cobra.Command{
		Use:  "send",
		RunE: cmdr.runIBCTransfer,
	}
	cmd.Flags().String(flagTo, "", "Address to send coins")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send")
	return cmd
}

type commander struct {
	cdc *wire.Codec
}

func (c commander) runIBCTransfer(cmd *cobra.Command, args []string) error {
	txBytes, err := c.buildTx()
	if err != nil {
		return err
	}

	res, err := client.BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
}

func buildMsg(from crypto.Address) (sdk.Msg, error) {
	amount := viper.GetString(flagAmount)
	coins, err := sdk.ParseCoins(amount)
	if err != nil {
		return nil, err
	}

	dest := viper.GetString(flagTo)
	bz, err := hex.DecodeString(dest)
	if err != nil {
		return nil, err
	}
	to := crypto.Address(bz)

	msg := ibc.IBCOutMsg{
		IBCTransfer: ibc.IBCTransfer{
			Destination: to,
			Coins:       coins,
		},
	}
}
