package order

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto"
	"testing"

	"github.com/ok-chain/okchain/x/token"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	abci "github.com/tendermint/tendermint/abci/types"
)

type MockDexApp struct {
	*mock.App

	keyOrders *sdk.KVStoreKey
	keyTrades *sdk.KVStoreKey

	keyToken     *sdk.KVStoreKey
	keyFreeze    *sdk.KVStoreKey
	keyLock      *sdk.KVStoreKey
	keyTokenPair *sdk.KVStoreKey
	keyFeeDetail *sdk.KVStoreKey

	bankKeeper  bank.Keeper
	orderKeeper Keeper
	tokenKeeper token.Keeper
}

type testAccount struct {
	addrKeys    *mock.AddrKeys
	baseAccount *auth.BaseAccount
}

type TestAccountList []*testAccount

// initialize the mock application for this module
func getMockDexApp(t *testing.T, numGenAccs int) (mockDexApp *MockDexApp, keeper Keeper,
	genAccs []auth.Account, addrs []sdk.AccAddress, pubKeys []crypto.PubKey, privKeys []crypto.PrivKey) {

	mapp := mock.NewApp()

	mockDexApp = &MockDexApp{
		App:       mapp,
		keyOrders: sdk.NewKVStoreKey(OrderStoreKey),
		keyTrades: sdk.NewKVStoreKey(TradeStoreKey),

		keyToken:     sdk.NewKVStoreKey("token"),
		keyFreeze:    sdk.NewKVStoreKey("freeze"),
		keyLock:      sdk.NewKVStoreKey("lock"),
		keyTokenPair: sdk.NewKVStoreKey("token_pair"),
		keyFeeDetail: sdk.NewKVStoreKey("fee_detail"),
	}

	mockDexApp.bankKeeper = bank.NewBaseKeeper(mockDexApp.AccountKeeper,
		mockDexApp.ParamsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace)

	mockDexApp.tokenKeeper = token.NewKeeper(
		mockDexApp.bankKeeper,
		mockDexApp.ParamsKeeper,
		mockDexApp.ParamsKeeper.Subspace(token.DefaultParamspace),
		mockDexApp.FeeCollectionKeeper,
		mockDexApp.keyToken,
		mockDexApp.keyFreeze,
		mockDexApp.keyLock,
		mockDexApp.keyTokenPair,
		mockDexApp.keyFeeDetail,
		mockDexApp.Cdc)

	mockDexApp.orderKeeper = NewKeeper(
		mockDexApp.tokenKeeper,
		mockDexApp.ParamsKeeper,
		mockDexApp.ParamsKeeper.Subspace(DefaultParamspace),
		mockDexApp.FeeCollectionKeeper,
		mockDexApp.keyOrders,
		mockDexApp.keyTrades,
		mockDexApp.Cdc)
	keeper = mockDexApp.orderKeeper

	RegisterCodec(mockDexApp.Cdc)

	mockDexApp.Router().AddRoute(RouterKey, NewHandler(mockDexApp.orderKeeper))
	mockDexApp.QueryRouter().AddRoute(QuerierRoute, NewQuerier(mockDexApp.orderKeeper))

	mockDexApp.SetEndBlocker(getEndBlocker(mockDexApp.orderKeeper))
	mockDexApp.SetInitChainer(getInitChainer(mockDexApp.App))

	intQuantity := int64(100)
	valTokens := token.ToUnit(intQuantity)
	coins := sdk.Coins{
		sdk.NewCoin("okb", valTokens),
		sdk.NewCoin("xxb", valTokens),
	}

	genAccs, addrs, pubKeys, privKeys = mock.CreateGenAccounts(numGenAccs, coins)

	// todo: checkTx in mock app
	mockDexApp.SetAnteHandler(nil)

	app := mockDexApp
	mockDexApp.MountStores(
		//app.keyOrders,
		//app.keyTrades,
		app.keyToken,
		app.keyTokenPair,
		app.keyFreeze,
		app.keyLock,
		app.keyFeeDetail,
	)

	require.NoError(t, mockDexApp.CompleteSetup(mockDexApp.keyOrders, mockDexApp.keyTrades))
	// TODO: set genesis
	mock.SetGenesis(mockDexApp.App, genAccs)
	return
}

func getEndBlocker(keeper Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		tags := EndBlocker(ctx, keeper)
		return abci.ResponseEndBlock{
			Tags: tags,
		}
	}
}

func getInitChainer(mapp *mock.App) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		return mapp.InitChainer(ctx, req)
	}
}
