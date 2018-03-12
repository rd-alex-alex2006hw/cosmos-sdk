package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tm "github.com/tendermint/tendermint/types"
)

type RegisterChainMsg struct {
}

type UpdateChainMsg struct {
	Header tm.Header
	Commit tm.Commit
}

type PacketCreateMsg struct {
	Packet
}

type PacketPostMsg struct {
	FromChainID     string
	FromChainHeight uint64
	Packet
	Proof *merkle.IAVLProof
}
