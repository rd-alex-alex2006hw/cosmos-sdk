package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/tendermint/go-crypto/keys"
	wire "github.com/tendermint/go-wire"

	"github.com/cosmos/cosmos-sdk/x/ibc"
)

func IBCRelayCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := commander{cdc, "ingress", "egress"}

	cmd := &cobra.Command{
		Use:  "relay",
		RunE: cmdr.runIBCTransfer,
	}
	cmd.Flags().String(flagTo, "", "Address to send coins")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send")
	return cmd
}

type commander struct {
	cdc          *wire.Codec
	keybase      keys.Keybase
	ingressStore string
	egressStore  string
}

func (c commander) runIBCRelay(cmd *cobra.Command, args []string) error {
	chain1 := viper.GetString(flagChain1)
	chain2 := viper.GetString(flagChain2)

	keybase, err := keys.GetKeyBase()
	if err != nil {
		return err
	}

	go loop(keybase, chain1, chain2)
	go loop(keybase, chain2, chain1)
}

// https://github.com/cosmos/cosmos-sdk/blob/master/client/helpers.go using specified address

func query(fromID string, key []byte, storeName string) (res []byte, err error) {
	orig := viper.GetSstring(client.FlagNode)
	viper.Set(client.FlagNode, fromID)
	res, err = client.Query(key, storeName)
	viper.Set(client.FlagNode, orig)
	return res, err
}

func broadcastTx(toID string, tx []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	orig := viper.GetSstring(client.FlagNode)
	viper.Set(client.FlagNode, toID)
	res, err := client.BroadcastTx(tx)
	viper.Set(client.FlagNode, orig)
	return res, err
}

func (c commander) refine(bz []byte) []byte {
	var transfer ibc.IBCTransfer
	if err = c.cdc.UnmarshalBinary(bz, &transfer); err != nil {
		panic(err)
	}
	msg := ibc.IBCInMsg{
		transfer,
	}
	res, err := buildTx(c.cdc, mssg)
	if err != nil {
		panic(err)
	}
	return res
}

func (c commander) loop(keybase keys.Keybase, fromID, toID string) {
	egressLengthKey, err := c.cdc.MarshalBinary(ibc.EgressKey(-1))
	if err != nil {
		panic(err)
	}

	ingressKey, err := c.cdc.MarshalBbinary(ibc.IngressKey{fromID})
	if err != nil {
		panic(err)
	}

	processed, err := c.query(to, ingressKey, c.ingressName)
	if err != nil {
		panic(err)
	}

OUTER:
	for {
		time.Sleep(time.Second)

		egressLength, err := c.query(fromNode, lengthKey, c.egressName)
		if err != nil {
			fmt.Printf("Error querying outgoing msg list length: '%s'\n", err)
			continue OUTER
		}

		for i := processed; i < egressLength; i++ {
			egressKey, err := c.query(from, EgressKey(i), c.egressName)
			if err != nil {
				panic(err)
			}

			bz, err := c.query(fromNode, egressKey, c.egressName)
			if err != nil {
				fmt.Printf("Error querying outgoing msg: '%s'\n", err)
				continue OUTER
			}

			_, err = c.broadcastTx(toNode, c.refine(bz))
			if err != nil {
				fmt.Printf("Error broadcasting incoming msg: '%s'\n", err)
				continue OUTER
			}

			fmt.Printf("Relayed msg: %d\n", i)
		}

		processed = egressLength
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
