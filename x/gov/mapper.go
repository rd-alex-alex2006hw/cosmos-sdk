package gov

import (
	"github.com/cosmos/cosmos-sdk/x/bank"
	wire "github.com/tendermint/go-wire"
)

type governanceMapper struct {
	ck bank.CoinKeeper

	// The (unexposed) key used to access the store from the Context.
	proposalStoreKey sdk.StoreKey

	// The (unexposed) key used to access the store from the Context.
	validatorInfoStoreKey sdk.StoreKey

	// The (unexposed) key used to access the store from the Context.
	votesStoreKey sdk.StoreKey

	// The (unexposed) key used to access the store from the Context.
	proposalProcessingQueueStoreKey sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec
}

// NewGovernanceMapper returns a mapper that uses go-wire to (binary) encode and decode gov types.
func NewGovernanceMapper(key sdk.StoreKey, ck bank.CoinKeeper) accountMapper {
	cdc := wire.NewCodec()
	return accountMapper{
		key: key,
		ck:  ck,
		cdc: cdc,
	}
}

// Returns the go-wire codec.  You may need to register interfaces
// and concrete types here, if your app's sdk.Account
// implementation includes interface fields.
// NOTE: It is not secure to expose the codec, so check out
// .Seal().
func (gm governanceMapper) WireCodec() *wire.Codec {
	return gm.cdc
}

func (gm governanceMapper) GetProposal(ctx sdk.Context, proposalId int64) sdk.Account {
	store := ctx.KVStore(am.proposalStoreKey)
	bz := store.Get(proposalId)
	if bz == nil {
		return nil
	}

	proposal := Proposal{}
	err := gm.cdc.UnmarshalBinary(bz, proposal)
	if err != nil {
		panic(err)
	}

	return acc
}

// Implements sdk.AccountMapper.
func (gm governanceMapper) SetProposal(ctx sdk.Context, proposal Proposal) {
	proposalId := proposal.ProposalId
	store := ctx.KVStore(am.proposalStoreKey)

	bz, err := gm.cdc.MarshalBinary(proposal)
	if err != nil {
		panic(err)
	}

	store.Set(proposalId, bz)
}

func (gm governanceMapper) GetValidatorInfo(ctx sdk.Context, proposalId int64, validatorAddr crypto.address) sdk.Account {
	store := ctx.KVStore(am.validatorInfoStoreKey)

	bz := store.Get(proposalId)
	if bz == nil {
		return nil
	}

	vote := Vote{}
	err := gm.cdc.UnmarshalBinary(bz, vote)
	if err != nil {
		panic(err)
	}

	return vote
}

// Implements sdk.AccountMapper.
func (gm governanceMapper) SetVote(ctx sdk.Context, vote Vote) {
	proposalId := proposal.ProposalId
	store := ctx.KVStore(am.votesStoreKey)

	bz, err := gm.cdc.MarshalBinary(vote)
	if err != nil {
		panic(err)
	}

	store.Set(proposalId, bz)
}

func (gm governanceMapper) GetVote(ctx sdk.Context, proposalId int64, voter crypto.address) sdk.Account {
	store := ctx.KVStore(am.votesStoreKey)
	bz := store.Get(proposalId)
	if bz == nil {
		return nil
	}

	vote := Vote{}
	err := gm.cdc.UnmarshalBinary(bz, vote)
	if err != nil {
		panic(err)
	}

	return vote
}

// Implements sdk.AccountMapper.
func (gm governanceMapper) SetVote(ctx sdk.Context, vote Vote) {
	proposalId := proposal.ProposalId
	store := ctx.KVStore(am.votesStoreKey)

	bz, err := gm.cdc.MarshalBinary(vote)
	if err != nil {
		panic(err)
	}

	store.Set(proposalId, bz)
}

type proposalQueueInfo struct {
	// begin <= elems < end
	begin int64
	end   int64
}

func (info proposalQueueInfo) validateBasic() error {
	if info.end < info.begin || info.begin < 0 || info.end < 0 {
		return errors.New("")
	}
	return nil
}

func (info proposalQueueInfo) isEmpty() bool {
	return begin == end
}

type proposalQueueElem int64

const proposalQueueInfoKey = int64(-1)

func (gm governanceMapper) getProposalInfo(store sdk.KVStore) proposalQueueInfo {
	bz := store.Get(proposalQueueInfoKey)
	info := proposalQueueInfo{}
	if err := gm.cdc.UnmarshalBinary(bz, &info); err != nil {
		panic(err)
	}
	if err = info.ValidateBasic(); err != nil {
		panic(err)
	}
	return info
}

func (gm governanceMapper) setProposalInfo(store sdk.KVStore, info proposalQueueInfo) {
	bz, err := gm.cdc.MarshalBinary(info)
	if err != nil {
		panic(err)
	}
	store.Set(proposalQueueInfoKey, bz)
}

func (gm governanceMapper) getProposalElem(store sdk.KVStore, index int64) int64 {
	return store.Get(index)
}

func (gm governanceMapper) setProposalElem(store sdk.KVStore, index int64, elem int64) {
	store.Set(index, elem)
}

func (gm governanceMapper) PeekProposalQueue(ctx sdk.Context) *int64 {
	store := ctx.KVStore(gm.proposalProcessingQueueStoreKey)

	info := gm.getProposalInfo(store)
	if info.isEmpty() {
		return nil
	}

	elem := gm.getProposalElem(store, info.begin)
	return &elem
}

func (gm governanceMapper) PushProposalQueue(ctx sdk.Context, proposalId int64) {
	store := ctx.KVStore(gm.proposalProcessingQueueStoreKey)

	info := getProposalInfo(store)
	setProposalElem(store, info.end, proposalId)

	info.end++
	gm.setProposalElem(store, info.end, proposalId)
}

func (gm governanceMapper) PopProposalQueue(ctx sdk.Context) {
	store := ctx.KVStore(gm.proposalProcessingQueueStoreKey)

	info := getProposalInfo(store)
	if info.isEmpty() {
		panic(errors.New(""))
	}

	store.Delete(info.begin)

	info.begin++
	gm.setProposalInfo(store, info)
}

// re-exporting github.com/cosmos/cosmos-sdk/x/bank/mapper.go

func (gm governanceMapper) SubtractCoins(ctx sdk.Context, addr crypto.Address, amt sdk.Coins) (sdk.Coins, sdk.Error) {
	return gm.ck.SubtractCoins(ctx, addr, amt)
}

func (gm governanceMapper) AddCoins(ctx sdk.Context, addr crypto.Address, amt sdk.Coins) (sdk.Coins, sdk.Error) {
	return gm.ck.AddCoins(ctx, addr, amt)
}
