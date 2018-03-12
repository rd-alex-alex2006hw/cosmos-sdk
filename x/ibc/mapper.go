package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type IBCMapper struct {
	ingressKey sdk.StoreKey
	egressKey  sdk.StoreKey
}

func (ibcm IBCMapper) GetSequenceNumber
