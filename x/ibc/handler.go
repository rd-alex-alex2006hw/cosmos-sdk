package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(ibcm sdk.StoreKey) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		return sdk.Result{}
	}
}
