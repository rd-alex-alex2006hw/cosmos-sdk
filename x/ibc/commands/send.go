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

func (c commander) buildTx() ([]byte, error) {
	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}

	name := viper.GetString(client.FlagName)
	info, err := keybase.Get(name)
	if err != nil {
		return nil, errors.Errorf("No key for: %s, name")
	}
	from := info.PubKey.Address()

	msg, err := buildMsg(from)
	if err != nil {
		return nil, err
	}

	bz := msg.GetSignBytes()
	buf := client.BufferStdin()
	prompt := fmt.Sprintf("Password to sign with '%s':", name)
	passphrase, err := client.GetPassword(prompt, buf)
	if err != nil {
		return nil, err
	}
	sig, pubkey, err := keybase.Sign(name, passphrase, bz)
	if err != nil {
		return nil, err
	}
	sigs := []sdk.StdSignature{{
		PubKey:    pubkey,
		Signature: sig,
		Sequence:  viper.GetInt64(flagSequence),
	}}

	tx := sdk.NewStdTx(msg, sigs)

	txBytes, err := c.cdc.MarshalBinary(tx)
	if err != nil {
		return nil, err
	}
	return txBytes, nil
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
		Destination: dest,
		Coins:       coins,
	}

	return msg, nil
}
