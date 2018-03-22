package eth

import (
	//sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/types"
)

// Wraps the RLP serialises representation of Ethereum transaction
type RawTxMsg struct {

	t types.Transaction
}