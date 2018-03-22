package eth

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// Wraps the RLP serialises representation of Ethereum transaction
type RawTxMsg struct {
	raw []byte
}

// enforce the msg type at compile time
var _ sdk.Msg = RawTxMsg{}

const EthRawMsgType = "ethraw"

// nolint
func (msg RawTxMsg) Type() string                            { return EthRawMsgType }
func (msg RawTxMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg RawTxMsg) GetSigners() []sdk.Address               { return []sdk.Address{} }
func (msg RawTxMsg) String() string {
	return fmt.Sprintf("RawTxMsg{%x}", msg.raw)
}

// Validate Basic is used to quickly disqualify obviously invalid messages
func (msg RawTxMsg) ValidateBasic() sdk.Error {
	return nil
}

// Get the bytes for the message signer to sign on
func (msg RawTxMsg) GetSignBytes() []byte {
	return msg.raw
}

func (msg RawTxMsg) DecodeRaw() (*ethtypes.Transaction, error) {
	return nil, nil
}
