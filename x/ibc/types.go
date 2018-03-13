package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	//	tm "github.com/tendermint/tendermint/types"
)

type IBCOutMsg struct {
	IBCTransfer
}

type IBCInMsg struct {
	IBCTransfer
}

type IBCTransfer struct {
	Destination sdk.Address
	Coins       sdk.Coins
}
