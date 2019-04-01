package app

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/okchain/okdex/x/backend"
	"github.com/okchain/okdex/x/gov"
	"github.com/okchain/okdex/x/order"
	"github.com/okchain/okdex/x/token"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	appName        = "dex"
	DefaultKeyPass = "12345678"
)

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.okdexcli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.okdexd")
)

type DexApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	keyMain          *types.KVStoreKey
	keyAccount       *types.KVStoreKey
	keyOrders        *types.KVStoreKey
	keyTrades        *types.KVStoreKey
	keyFeeCollection *types.KVStoreKey
	keyParams        *types.KVStoreKey
	keyToken         *types.KVStoreKey
	keyFreeze        *types.KVStoreKey
	keyLock          *types.KVStoreKey
	keyTokenPair     *types.KVStoreKey
	tkeyParams       *types.TransientStoreKey
	keyDistr         *types.KVStoreKey
	tkeyDistr        *types.TransientStoreKey
	keyStaking       *types.KVStoreKey
	tkeyStaking      *types.TransientStoreKey
	keySlashing      *types.KVStoreKey
	keyGov           *types.KVStoreKey
	keyMint          *types.KVStoreKey

	accountKeeper       auth.AccountKeeper
	bankKeeper          bank.Keeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	paramsKeeper        params.Keeper
	orderKeeper         order.Keeper
	backendKeeper       backend.Keeper
	tokenKeeper         token.Keeper
	distrKeeper         distr.Keeper
	stakingKeeper       staking.Keeper
	slashingKeeper      slashing.Keeper
	mintKeeper          mint.Keeper
	govKeeper           gov.Keeper
}

// DecentralizedExchangeApp is a constructor function for DEXApp
func DecentralizedExchangeApp(logger log.Logger, db dbm.DB) *DexApp {

	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	// Here you initialize your application with the store keys it requires
	var app = &DexApp{
		BaseApp: bApp,
		cdc:     cdc,

		keyMain:          sdk.NewKVStoreKey("main"),
		keyAccount:       sdk.NewKVStoreKey(auth.StoreKey),
		keyOrders:        sdk.NewKVStoreKey(order.OrderStoreKey),
		keyTrades:        sdk.NewKVStoreKey(order.TradeStoreKey),
		keyFeeCollection: sdk.NewKVStoreKey("fee_collection"),
		keyParams:        sdk.NewKVStoreKey("params"),
		keyToken:         sdk.NewKVStoreKey("token"),
		keyFreeze:        sdk.NewKVStoreKey("freeze"),
		keyLock:          sdk.NewKVStoreKey("lock"),
		keyTokenPair:     sdk.NewKVStoreKey("token_pair"),
		tkeyParams:       sdk.NewTransientStoreKey("transient_params"),
		keyStaking:       sdk.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:      sdk.NewTransientStoreKey(staking.TStoreKey),
		keyMint:          sdk.NewKVStoreKey(mint.StoreKey),
		keyDistr:         sdk.NewKVStoreKey(distr.StoreKey),
		tkeyDistr:        sdk.NewTransientStoreKey(distr.TStoreKey),
		keySlashing:      sdk.NewKVStoreKey(slashing.StoreKey),
		keyGov:           sdk.NewKVStoreKey(gov.StoreKey),
	}

	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tkeyParams)

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,
		app.paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)
	stakingKeeper := staking.NewKeeper(
		app.cdc,
		app.keyStaking, app.tkeyStaking,
		app.bankKeeper, app.paramsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)
	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		app.keyDistr,
		app.paramsKeeper.Subspace(distr.DefaultParamspace),
		app.bankKeeper, &stakingKeeper, app.feeCollectionKeeper,
		distr.DefaultCodespace,
	)

	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		app.keySlashing,
		&stakingKeeper, app.paramsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
	)
	app.govKeeper = gov.NewKeeper(
		app.cdc,
		app.keyGov,
		app.paramsKeeper,
		&app.tokenKeeper,
		app.paramsKeeper.Subspace(gov.DefaultParamspace),
		app.bankKeeper,
		&stakingKeeper,
		gov.DefaultCodespace,
	)
	app.mintKeeper = mint.NewKeeper(app.cdc, app.keyMint,
		app.paramsKeeper.Subspace(mint.DefaultParamspace),
		&stakingKeeper, app.feeCollectionKeeper,
	)
	// register the staking hooks
	// NOTE: The stakingKeeper above is passed by reference, so that it can be
	// modified like below:
	app.stakingKeeper = stakingKeeper

	// The FeeCollectionKeeper collects transaction fees and renders them to the fee distribution module
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(cdc, app.keyFeeCollection)

	// The NameserviceKeeper is the Keeper from the module for this tutorial
	// It handles interactions with the namestore

	app.tokenKeeper = token.NewKeeper(
		app.bankKeeper,
		app.paramsKeeper,
		app.paramsKeeper.Subspace(token.DefaultParamspace),
		app.feeCollectionKeeper,
		app.keyToken,
		app.keyFreeze,
		app.keyLock,
		app.keyTokenPair,
		app.cdc,
	)
	app.orderKeeper = order.NewKeeper(
		app.tokenKeeper,
		app.paramsKeeper,
		app.paramsKeeper.Subspace(order.DefaultParamspace),
		app.feeCollectionKeeper,
		app.keyOrders,
		app.keyTrades,
		app.cdc,
	)
	app.backendKeeper = backend.NewKeeper(
		app.orderKeeper,
		app.cdc,
	)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))

	app.paramsKeeper.RegisterParamSet(gov.DefaultParamspace, &gov.GovParams{}).
		RegisterParamSet(auth.DefaultParamspace, &auth.Params{}).
		// RegisterParamSet(bank.DefaultParamspace, &bank.ParamStoreKeySendEnabled).
		RegisterParamSet(staking.DefaultParamspace, &staking.Params{}).
		// RegisterParamSet(distr.DefaultParamspace, &distr.Params{}).
		RegisterParamSet(slashing.DefaultParamspace, &slashing.Params{}).
		// RegisterParamSet(mint.DefaultParamspace, &mint.Params{}).
		RegisterParamSet(order.DefaultParamspace, &order.Params{}).
		RegisterParamSet(token.DefaultParamspace, &token.Params{})

	// The app.Router is the main transaction router where each module registers its routes
	// Register the bank and nameservice routes here
	app.Router().
		AddRoute(bank.RouterKey, bank.NewHandler(app.bankKeeper)).
		AddRoute(order.RouterKey, order.NewHandler(app.orderKeeper)).
		AddRoute(token.RouterKey, token.NewHandler(app.tokenKeeper)).
		AddRoute(distr.RouterKey, distr.NewHandler(app.distrKeeper)).
		AddRoute(staking.RouterKey, staking.NewHandler(app.stakingKeeper)).
		AddRoute(gov.RouterKey, gov.NewHandler(app.govKeeper))

	// The app.QueryRouter is the main query router where each module registers its routes
	app.QueryRouter().
		AddRoute(order.QuerierRoute, order.NewQuerier(app.orderKeeper)).
		AddRoute(token.QuerierRoute, token.NewQuerier(app.tokenKeeper)).
		AddRoute(auth.QuerierRoute, auth.NewQuerier(app.accountKeeper)).
		AddRoute(staking.QuerierRoute, staking.NewQuerier(app.stakingKeeper, app.cdc)).
		AddRoute(distr.QuerierRoute, distr.NewQuerier(app.distrKeeper)).
		AddRoute(gov.QuerierRoute, gov.NewQuerier(app.govKeeper)).
		AddRoute(slashing.QuerierRoute, slashing.NewQuerier(app.slashingKeeper, app.cdc)).
		AddRoute(backend.QuerierRoute, backend.NewQuerier(app.backendKeeper))

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.initChainer)

	// application updates every end block
	app.SetEndBlocker(app.endBlocker)

	app.MountStores(
		app.keyMain,
		app.keyAccount,
		app.keyOrders,
		app.keyTrades,
		app.keyToken,
		app.keyTokenPair,
		app.keyFreeze,
		app.keyLock,
		app.keyFeeCollection,
		app.keyParams,
		app.tkeyParams,
		app.keyStaking,
		app.keyMint,
		app.keyDistr,
		app.keySlashing,
		app.keyGov,
		app.tkeyDistr,
	)

	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

func (app *DexApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	var genesisState GenesisState
	err := app.cdc.UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err)
	}

	validators := app.initFromGenesisState(ctx, genesisState)

	for _, token := range genesisState.Token.Info {
		if token.Name != "" && token.Symbol != "" && token.TotalSupply > 0 && token.Owner != nil {
			err = app.issueOKB(ctx, token)
			if err != nil {
				panic(err)
			}
		}
	}
	// sanity check
	//if len(req.Validators) > 0 {
	//	if len(req.Validators) != len(validators) {
	//		panic(fmt.Errorf("len(RequestInitChain.Validators) != len(validators) (%d != %d)",
	//			len(req.Validators), len(validators)))
	//	}
	//	sort.Sort(abci.ValidatorUpdates(req.Validators))
	//	sort.Sort(abci.ValidatorUpdates(validators))
	//	for i, val := range validators {
	//		if !val.Equal(req.Validators[i]) {
	//			panic(fmt.Errorf("validators[%d] != req.Validators[%d] ", i, i))
	//		}
	//	}
	//}

	// assert runtime invariants
	// app.assertRuntimeInvariants()

	return abci.ResponseInitChain{
		Validators: validators,
	}
}

func (app *DexApp) issueOKB(ctx sdk.Context, okb token.Token) error {
	coins := app.tokenKeeper.GetCoins(ctx, okb.Owner)
	if !strings.Contains(coins.String(), okb.Symbol) {
		coins = append(coins, sdk.NewCoin(okb.Symbol, token.ToUnit(okb.TotalSupply)))
		sort.Sort(coins)

		err := app.tokenKeeper.SetCoins(ctx, okb.Owner, coins)
		if err != nil {
			return err
		}
	}

	app.tokenKeeper.NewToken(ctx, okb)
	return nil
}

func (app *DexApp) endBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	tags := gov.EndBlocker(ctx, app.govKeeper)
	tags.AppendTags(order.EndBlocker(ctx, app.orderKeeper))
	tags.AppendTags(backend.EndBlocker(ctx, app.backendKeeper, app.orderKeeper))

	return abci.ResponseEndBlock{
		Tags: tags,
	}
}

// ExportAppStateAndValidators does the things
func (app *DexApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})
	// iterate to get the accounts
	accounts := []GenesisAccount{}
	appendAccount := func(acc auth.Account) (stop bool) {
		account := NewGenesisAccountI(acc)
		accounts = append(accounts, account)
		return false
	}
	app.accountKeeper.IterateAccounts(ctx, appendAccount)

	genState := GenesisState{
		Accounts: accounts,
		AuthData: auth.DefaultGenesisState(),
		BankData: bank.DefaultGenesisState(),
		GovData:  gov.DefaultGenesisState(),
	}

	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	return appState, validators, err
}

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	order.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	token.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

// initialize store from a genesis state
func (app *DexApp) initFromGenesisState(ctx sdk.Context, genesisState GenesisState) []abci.ValidatorUpdate {
	genesisState.Sanitize()

	// load the accounts
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		acc = app.accountKeeper.NewAccount(ctx, acc) // set account number
		app.accountKeeper.SetAccount(ctx, acc)
	}

	// initialize distribution (must happen before staking)
	distr.InitGenesis(ctx, app.distrKeeper, genesisState.DistrData)

	// load the initial staking information
	validators, err := staking.InitGenesis(ctx, app.stakingKeeper, genesisState.StakingData)
	if err != nil {
		panic(err) // TODO find a way to do this w/o panics
	}

	// initialize module-specific stores
	auth.InitGenesis(ctx, app.accountKeeper, app.feeCollectionKeeper, genesisState.AuthData)
	bank.InitGenesis(ctx, app.bankKeeper, genesisState.BankData)
	slashing.InitGenesis(ctx, app.slashingKeeper, genesisState.SlashingData, genesisState.StakingData.Validators.ToSDKValidators())
	gov.InitGenesis(ctx, app.govKeeper, genesisState.GovData)
	mint.InitGenesis(ctx, app.mintKeeper, genesisState.MintData)
	order.InitGenesis(ctx, app.orderKeeper, genesisState.Order)
	token.InitGenesis(ctx, app.tokenKeeper, genesisState.Token)

	// validate genesis state
	// if err := ValidateGenesisState(genesisState); err != nil {
	// 	panic(err) // TODO find a way to do this w/o panics
	// }

	if len(genesisState.GenTxs) > 0 {
		for _, genTx := range genesisState.GenTxs {
			var tx auth.StdTx
			err = app.cdc.UnmarshalJSON(genTx, &tx)
			if err != nil {
				panic(err)
			}
			bz := app.cdc.MustMarshalBinaryLengthPrefixed(tx)
			res := app.BaseApp.DeliverTx(bz)
			if !res.IsOK() {
				panic(res.Log)
			}
		}

		validators = app.stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}
	return validators
}
