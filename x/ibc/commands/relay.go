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
	flagChain1 = "chain1"
	flagChain2 = "chain2"
)

func IBCRelayCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := commander{cdc}

	cmd := &cobra.Command{
		Use:  "relay",
		RunE: cmdr.runIBCTransfer,
	}
	cmd.Flags().String(flagTo, "", "Address to send coins")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send")
	return cmd
}

type commander struct {
	cdc       *wire.Codec
	storeName string
}

func (c commander) runIBCRelay(cmd *cobra.Command, args []string) error {
	chain1 := viper.GetString(flagChain1)
	chain2 := viper.GetString(flagChain2)

	keybase, err := keys.GetKeyBase()
	if err != nil {
		return err
	}

	node1 := rpcclient.NewHTTP(chain1, "/websocket")
	node2 := rpcclient.NewHTTP(chain2, "/websocket")

	go loop(keybase, node1, node2)
	go loop(keybase, node2, node1)
}

// https://github.com/cosmos/cosmos-sdk/blob/master/client/helpers.go using specified address

func (c commander) loop(keybase keys.Keybase, from, to rpcclient.Client) {
	key := -1

	nextSeq := 0

	for {
		time.Sleep(time.Second)

		key := nextSeq
		res, err := c.query(from, key, c.storeName)
	}
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
