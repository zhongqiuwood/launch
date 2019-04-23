package app

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/ok-chain/okchain/util"
	"github.com/ok-chain/okchain/x/perf"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// abci "github.com/tendermint/tendermint/abci/types"
// Application is an interface that enables any finite, deterministic state machine
// to be driven by a blockchain-based replication engine via the ABCI.
// All methods take a RequestXxx argument and return a ResponseXxx argument,
// except CheckTx/DeliverTx, which take `tx []byte`, and `Commit`, which takes nothing.
//type Application interface {
//    // Info/Query Connection
//    Info(RequestInfo) ResponseInfo                // Return application info
//    SetOption(RequestSetOption) ResponseSetOption // Set application option
//    Query(RequestQuery) ResponseQuery             // Query for state
//
//    // Mempool Connection
//    CheckTx(tx []byte) ResponseCheckTx // Validate a tx for the mempool
//
//    // Consensus Connection
//    InitChain(RequestInitChain) ResponseInitChain    // Initialize blockchain with validators and other info from TendermintCore
//    BeginBlock(RequestBeginBlock) ResponseBeginBlock // Signals the beginning of a block
//    DeliverTx(tx []byte) ResponseDeliverTx           // Deliver a tx for full processing
//    EndBlock(RequestEndBlock) ResponseEndBlock       // Signals the end of a block, returns changes to the validator set
//    Commit() ResponseCommit                          // Commit the state and return the application Merkle root hash
//}

func (app *DexApp) log(format string, a ...interface{}) {
	format = fmt.Sprintf("[%s]%s", util.GoId, format)
	app.Logger().Info(fmt.Sprintf(format, a...))
}

func (app *DexApp) InitChain(req abci.RequestInitChain) (res abci.ResponseInitChain) {

	app.log("[ABCI interface] ---> InitChain")
	return app.BaseApp.InitChain(req)
}

// CheckTx implements the ABCI interface. It runs the "basic checks" to see
// whether or not a transaction can possibly be executed, first decoding, then
// the ante handler (which checks signatures/fees/ValidateBasic), then finally
// the route match to see whether a handler exists.
//
// NOTE:CheckTx does not run the actual Msg handler function(s).

func (app *DexApp) CheckTx(txBytes []byte) (res abci.ResponseCheckTx) {

	//app.log("===============================")
	//app.log("[ABCI interface] ---> CheckTx in")
	//defer app.log("[ABCI interface] ---> CheckTx out")

	//time.Sleep(5 * time.Second)

	return app.BaseApp.CheckTx(txBytes)
}

// ===================================
// ===================================
// Consensus Connection

func (app *DexApp) BeginBlock(req abci.RequestBeginBlock) (res abci.ResponseBeginBlock) {

	app.log("[ABCI interface][%d] ---> BeginBlock in", app.LastBlockHeight())
	defer app.log("[ABCI interface][%d]  ---> BeginBlock out", app.LastBlockHeight())

	seq := perf.GetPerf().OnAppBeginBlockEnter(app.LastBlockHeight() + 1)
	defer perf.GetPerf().OnAppBeginBlockExit(app.LastBlockHeight() + 1, seq)

	return app.BaseApp.BeginBlock(req)
}

// sync txBytes to backend module
func (app *DexApp) syncTx(txBytes []byte) {
	if tx, err := auth.DefaultTxDecoder(app.cdc)(txBytes); err == nil {
		if stdTx, ok := tx.(auth.StdTx); ok {
			txHash := fmt.Sprintf("%X", tmhash.Sum(txBytes))
			app.log("[Sync Tx(%s) to backend module]", txHash)
			ctx := app.GetState(baseapp.RunTxModeDeliver).Context()
			app.backendKeeper.ProduceTxEvent(ctx, &stdTx, txHash, ctx.BlockHeader().Time.Unix())
		}
	}
}

func (app *DexApp) DeliverTx(txBytes []byte) (res abci.ResponseDeliverTx) {
	app.log("[ABCI interface] ---> DeliverTx in")
	defer app.log("[ABCI interface] ---> DeliverTx out")

	response := app.BaseApp.DeliverTx(txBytes)
	if response.IsOK() {
		app.syncTx(txBytes)
	}
	return response
}

// EndBlock implements the ABCI interface.
func (app *DexApp) EndBlock(req abci.RequestEndBlock) (res abci.ResponseEndBlock) {
	app.log("[ABCI interface] ---> EndBlock in")
	defer app.log("[ABCI interface] ---> EndBlock out")

	seq := perf.GetPerf().OnAppEndBlockEnter(app.LastBlockHeight() + 1)
	defer perf.GetPerf().OnAppEndBlockExit(app.LastBlockHeight() + 1, seq)

	return app.BaseApp.EndBlock(req)
}

// Commit implements the ABCI interface.
func (app *DexApp) Commit() abci.ResponseCommit {

	seq := perf.GetPerf().OnCommitEnter(app.LastBlockHeight() + 1)
	defer perf.GetPerf().OnCommitExit(app.LastBlockHeight() + 1, seq, app.Logger())

	app.log("[ABCI interface] ---> Commit in")
	res := app.BaseApp.Commit()
	defer app.log("[ABCI interface] ---> Commit out")

	return res
}


// Consensus Connection
// ===================================
// ===================================
