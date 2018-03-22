package app

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/examples/ethermint/types"
	"github.com/cosmos/cosmos-sdk/examples/ethermint/x/cool"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"

	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

// Construct some global addrs and txs for tests.
var (
	chainID = "" // TODO

	priv1 = crypto.GenPrivKeyEd25519()
	addr1 = priv1.PubKey().Address()
	addr2 = crypto.GenPrivKeyEd25519().PubKey().Address()
	coins = sdk.Coins{{"foocoin", 10}}
	fee   = sdk.StdFee{
		sdk.Coins{{"foocoin", 0}},
		0,
	}

	sendMsg = bank.SendMsg{
		Inputs:  []bank.Input{bank.NewInput(addr1, coins)},
		Outputs: []bank.Output{bank.NewOutput(addr2, coins)},
	}

	quizMsg1 = cool.QuizMsg{
		Sender:     addr1,
		CoolAnswer: "icecold",
	}

	quizMsg2 = cool.QuizMsg{
		Sender:     addr1,
		CoolAnswer: "badvibesonly",
	}

	setTrendMsg1 = cool.SetTrendMsg{
		Sender: addr1,
		Cool:   "icecold",
	}

	setTrendMsg2 = cool.SetTrendMsg{
		Sender: addr1,
		Cool:   "badvibesonly",
	}

	setTrendMsg3 = cool.SetTrendMsg{
		Sender: addr1,
		Cool:   "warmandkind",
	}
)

func newEthermintApp() *EthermintApp {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	return NewEthermintApp(logger, db)
}

//_______________________________________________________________________

func TestMsgs(t *testing.T) {
	eapp := newEthermintApp()

	msgs := []struct {
		msg sdk.Msg
	}{
		{sendMsg},
		{quizMsg1},
		{setTrendMsg1},
	}

	sequences := []int64{0}
	for i, m := range msgs {
		sig := priv1.Sign(sdk.StdSignBytes(chainID, sequences, fee, m.msg))
		tx := sdk.NewStdTx(m.msg, fee, []sdk.StdSignature{{
			PubKey:    priv1.PubKey(),
			Signature: sig,
		}})

		// just marshal/unmarshal!
		cdc := MakeCodec()
		txBytes, err := cdc.MarshalBinary(tx)
		require.NoError(t, err, "i: %v", i)

		// Run a Check
		cres := eapp.CheckTx(txBytes)
		assert.Equal(t, sdk.CodeUnknownAddress,
			sdk.CodeType(cres.Code), "i: %v, log: %v", i, cres.Log)

		// Simulate a Block
		eapp.BeginBlock(abci.RequestBeginBlock{})
		dres := eapp.DeliverTx(txBytes)
		assert.Equal(t, sdk.CodeUnknownAddress,
			sdk.CodeType(dres.Code), "i: %v, log: %v", i, dres.Log)
	}
}

func TestGenesis(t *testing.T) {
	eapp := newEthermintApp()

	// Construct some genesis bytes to reflect ethermint/types/AppAccount
	pk := crypto.GenPrivKeyEd25519().PubKey()
	addr := pk.Address()
	coins, err := sdk.ParseCoins("77foocoin,99barcoin")
	require.Nil(t, err)
	baseAcc := auth.BaseAccount{
		Address: addr,
		Coins:   coins,
	}
	acc := &types.AppAccount{baseAcc, "foobart"}

	genesisState := types.GenesisState{
		Accounts: []*types.GenesisAccount{
			types.NewGenesisAccount(acc),
		},
	}
	stateBytes, err := json.MarshalIndent(genesisState, "", "\t")

	vals := []abci.Validator{}
	eapp.InitChain(abci.RequestInitChain{vals, stateBytes})
	eapp.Commit()

	// A checkTx context
	ctx := eapp.BaseApp.NewContext(true, abci.Header{})
	res1 := eapp.accountMapper.GetAccount(ctx, baseAcc.Address)
	assert.Equal(t, acc, res1)
	/*
		// reload app and ensure the account is still there
		bapp = NewBasecoinApp(logger, db)
		ctx = bapp.BaseApp.NewContext(true, abci.Header{})
		res1 = bapp.accountMapper.GetAccount(ctx, baseAcc.Address)
		assert.Equal(t, acc, res1)
	*/
}

func TestSendMsgWithAccounts(t *testing.T) {
	eapp := newEthermintApp()

	// Construct some genesis bytes to reflect basecoin/types/AppAccount
	// Give 77 foocoin to the first key
	coins, err := sdk.ParseCoins("77foocoin")
	require.Nil(t, err)
	baseAcc := auth.BaseAccount{
		Address: addr1,
		Coins:   coins,
	}
	acc1 := &types.AppAccount{baseAcc, "foobart"}

	// Construct genesis state
	genesisState := types.GenesisState{
		Accounts: []*types.GenesisAccount{
			types.NewGenesisAccount(acc1),
		},
	}
	stateBytes, err := json.MarshalIndent(genesisState, "", "\t")
	require.Nil(t, err)

	// Initialize the chain
	vals := []abci.Validator{}
	eapp.InitChain(abci.RequestInitChain{vals, stateBytes})
	eapp.Commit()

	// A checkTx context (true)
	ctxCheck := eapp.BaseApp.NewContext(true, abci.Header{})
	res1 := eapp.accountMapper.GetAccount(ctxCheck, addr1)
	assert.Equal(t, acc1, res1)

	// Sign the tx
	sequences := []int64{0}
	sig := priv1.Sign(sdk.StdSignBytes(chainID, sequences, fee, sendMsg))
	tx := sdk.NewStdTx(sendMsg, fee, []sdk.StdSignature{{
		PubKey:    priv1.PubKey(),
		Signature: sig,
	}})

	// Run a Check
	res := eapp.Check(tx)
	assert.Equal(t, sdk.CodeOK, res.Code, res.Log)

	// Simulate a Block
	eapp.BeginBlock(abci.RequestBeginBlock{})
	res = eapp.Deliver(tx)
	assert.Equal(t, sdk.CodeOK, res.Code, res.Log)

	// Check balances
	ctxDeliver := eapp.BaseApp.NewContext(false, abci.Header{})
	res2 := eapp.accountMapper.GetAccount(ctxDeliver, addr1)
	res3 := eapp.accountMapper.GetAccount(ctxDeliver, addr2)
	assert.Equal(t, fmt.Sprintf("%v", res2.GetCoins()), "67foocoin")
	assert.Equal(t, fmt.Sprintf("%v", res3.GetCoins()), "10foocoin")

	// Delivering again should cause replay error
	res = eapp.Deliver(tx)
	assert.Equal(t, sdk.CodeInvalidSequence, res.Code, res.Log)

	// bumping the txnonce number without resigning should be an auth error
	tx.Signatures[0].Sequence = 1
	res = eapp.Deliver(tx)
	assert.Equal(t, sdk.CodeUnauthorized, res.Code, res.Log)

	// resigning the tx with the bumped sequence should work
	sequences = []int64{1}
	sig = priv1.Sign(sdk.StdSignBytes(chainID, sequences, fee, tx.Msg))
	tx.Signatures[0].Signature = sig
	res = eapp.Deliver(tx)
	assert.Equal(t, sdk.CodeOK, res.Code, res.Log)
}

func TestQuizMsg(t *testing.T) {
	eapp := newEthermintApp()

	// Construct genesis state
	// Construct some genesis bytes to reflect basecoin/types/AppAccount
	coins := sdk.Coins{}
	baseAcc := auth.BaseAccount{
		Address: addr1,
		Coins:   coins,
	}
	acc1 := &types.AppAccount{baseAcc, "foobart"}

	// Construct genesis state
	genesisState := types.GenesisState{
		Accounts: []*types.GenesisAccount{
			types.NewGenesisAccount(acc1),
		},
	}
	stateBytes, err := json.MarshalIndent(genesisState, "", "\t")
	require.Nil(t, err)

	// Initialize the chain (nil)
	vals := []abci.Validator{}
	eapp.InitChain(abci.RequestInitChain{vals, stateBytes})
	eapp.Commit()

	// A checkTx context (true)
	ctxCheck := eapp.BaseApp.NewContext(true, abci.Header{})
	res1 := eapp.accountMapper.GetAccount(ctxCheck, addr1)
	assert.Equal(t, acc1, res1)

	// Set the trend, submit a really cool quiz and check for reward
	SignCheckDeliver(t, eapp, setTrendMsg1, 0, true)
	SignCheckDeliver(t, eapp, quizMsg1, 1, true)
	CheckBalance(t, eapp, "69icecold")
	SignCheckDeliver(t, eapp, quizMsg2, 2, true) // result without reward
	CheckBalance(t, eapp, "69icecold")
	SignCheckDeliver(t, eapp, quizMsg1, 3, true)
	CheckBalance(t, eapp, "138icecold")
	SignCheckDeliver(t, eapp, setTrendMsg2, 4, true) // reset the trend
	SignCheckDeliver(t, eapp, quizMsg1, 5, true)     // the same answer will nolonger do!
	CheckBalance(t, eapp, "138icecold")
	SignCheckDeliver(t, eapp, quizMsg2, 6, true) // earlier answer now relavent again
	CheckBalance(t, eapp, "69badvibesonly,138icecold")
	SignCheckDeliver(t, eapp, setTrendMsg3, 7, false) // expect to fail to set the trend to something which is not cool

}

func TestHandler(t *testing.T) {
	eapp := newEthermintApp()

	sourceChain := "source-chain"
	destChain := "dest-chain"

	vals := []abci.Validator{}
	baseAcc := auth.BaseAccount{
		Address: addr1,
		Coins:   coins,
	}
	acc1 := &types.AppAccount{baseAcc, "foobart"}
	genesisState := types.GenesisState{
		Accounts: []*types.GenesisAccount{
			types.NewGenesisAccount(acc1),
		},
	}
	stateBytes, err := json.MarshalIndent(genesisState, "", "\t")
	require.Nil(t, err)
	eapp.InitChain(abci.RequestInitChain{vals, stateBytes})
	eapp.Commit()

	// A checkTx context (true)
	ctxCheck := eapp.BaseApp.NewContext(true, abci.Header{})
	res1 := eapp.accountMapper.GetAccount(ctxCheck, addr1)
	assert.Equal(t, acc1, res1)

	packet := ibc.IBCPacket{
		SrcAddr:   addr1,
		DestAddr:  addr1,
		Coins:     coins,
		SrcChain:  sourceChain,
		DestChain: destChain,
	}

	transferMsg := ibc.IBCTransferMsg{
		IBCPacket: packet,
	}

	receiveMsg := ibc.IBCReceiveMsg{
		IBCPacket: packet,
		Relayer:   addr1,
		Sequence:  0,
	}

	SignCheckDeliver(t, eapp, transferMsg, 0, true)
	CheckBalance(t, eapp, "")
	SignCheckDeliver(t, eapp, transferMsg, 1, false)
	SignCheckDeliver(t, eapp, receiveMsg, 2, true)
	CheckBalance(t, eapp, "10foocoin")
	SignCheckDeliver(t, eapp, receiveMsg, 3, false)
}

func SignCheckDeliver(t *testing.T, eapp *EthermintApp, msg sdk.Msg, seq int64, expPass bool) {

	// Sign the tx
	tx := sdk.NewStdTx(msg, fee, []sdk.StdSignature{{
		PubKey:    priv1.PubKey(),
		Signature: priv1.Sign(sdk.StdSignBytes(chainID, []int64{seq}, fee, msg)),
		Sequence:  seq,
	}})

	// Run a Check
	res := eapp.Check(tx)
	if expPass {
		require.Equal(t, sdk.CodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.CodeOK, res.Code, res.Log)
	}

	// Simulate a Block
	eapp.BeginBlock(abci.RequestBeginBlock{})
	res = eapp.Deliver(tx)
	if expPass {
		require.Equal(t, sdk.CodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.CodeOK, res.Code, res.Log)
	}
}

func CheckBalance(t *testing.T, eapp *EthermintApp, balExpected string) {
	ctxDeliver := eapp.BaseApp.NewContext(false, abci.Header{})
	res2 := eapp.accountMapper.GetAccount(ctxDeliver, addr1)
	assert.Equal(t, balExpected, fmt.Sprintf("%v", res2.GetCoins()))
}
