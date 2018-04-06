package simplestake

import (
	"fmt"

	"testing"

	"github.com/stretchr/testify/assert"

	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	authKey := sdk.NewKVStoreKey("authkey")
	capKey := sdk.NewKVStoreKey("capkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(capKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, authKey, capKey
}

func TestKeeperGetSet(t *testing.T) {
	ms, _, capKey := setupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, nil, sdk.CodespaceRoot)
	stakeKeeper := NewKeeper(capKey, bank.NewCoinKeeper(nil))
	addr := sdk.Address([]byte("some-address"))

	bi := stakeKeeper.getBondInfo(ctx, addr)
	assert.Equal(t, bi, bondInfo{})

	privKey := crypto.GenPrivKeyEd25519()

	bi = bondInfo{
		PubKey: privKey.PubKey(),
		Power:  int64(10),
	}
	fmt.Printf("Pubkey: %v\n", privKey.PubKey())
	stakeKeeper.setBondInfo(ctx, addr, bi)

	savedBi := stakeKeeper.getBondInfo(ctx, addr)
	assert.NotNil(t, savedBi)
	fmt.Printf("Bond Info: %v\n", savedBi)
	assert.Equal(t, int64(10), savedBi.Power)
}

func TestBonding(t *testing.T) {
	ms, authKey, capKey := setupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, nil, sdk.CodespaceRoot)

	accountMapper := auth.NewAccountMapper(authKey, &auth.BaseAccount{})
	coinKeeper := bank.NewCoinKeeper(accountMapper)
	stakeKeeper := NewKeeper(capKey, coinKeeper)
	addr := sdk.Address([]byte("some-address"))
	privKey := crypto.GenPrivKeyEd25519()
	pubKey := privKey.PubKey()

	_, _, err := stakeKeeper.unbondWithoutCoins(ctx, addr)
	assert.Equal(t, err, ErrInvalidUnbond(sdk.CodespaceRoot))

	_, err = stakeKeeper.bondWithoutCoins(ctx, addr, pubKey, sdk.Coin{"steak", 10})
	assert.Nil(t, err)

	power, err := stakeKeeper.bondWithoutCoins(ctx, addr, pubKey, sdk.Coin{"steak", 10})
	assert.Equal(t, int64(20), power)

	pk, _, err := stakeKeeper.unbondWithoutCoins(ctx, addr)
	assert.Nil(t, err)
	assert.Equal(t, pubKey, pk)

	_, _, err = stakeKeeper.unbondWithoutCoins(ctx, addr)
	assert.Equal(t, err, ErrInvalidUnbond(sdk.CodespaceRoot))
}
