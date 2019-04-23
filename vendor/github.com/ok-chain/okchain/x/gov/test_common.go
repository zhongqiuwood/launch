package gov

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ok-chain/okchain/x/token"
	dbm "github.com/tendermint/tendermint/libs/db"
)

// initialize the mock application for this module
func getMockApp(t *testing.T, numGenAccs int, genState GenesisState, genAccs []auth.Account) (
	mapp *mock.App, keeper Keeper, sk staking.Keeper, addrs []sdk.AccAddress,
	pubKeys []crypto.PubKey, privKeys []crypto.PrivKey) {

	mapp = mock.NewApp()

	staking.RegisterCodec(mapp.Cdc)
	RegisterCodec(mapp.Cdc)

	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keyGov := sdk.NewKVStoreKey(StoreKey)
	keyToken := sdk.NewKVStoreKey("token")
	keyFreeze := sdk.NewKVStoreKey("freeze")
	keyLock := sdk.NewKVStoreKey("lock")
	keyTokenPair := sdk.NewKVStoreKey("token_pair")
	keyFeeDetail := sdk.NewKVStoreKey("fee_detail")

	pk := mapp.ParamsKeeper
	ck := bank.NewBaseKeeper(mapp.AccountKeeper, mapp.ParamsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace)
	sk = staking.NewKeeper(mapp.Cdc, keyStaking, tkeyStaking, ck, pk.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)
	tk := token.NewKeeper(ck,
		mapp.ParamsKeeper,
		mapp.ParamsKeeper.Subspace(token.DefaultParamspace),
		mapp.FeeCollectionKeeper,
		keyToken,
		keyFreeze,
		keyLock,
		keyTokenPair,
		keyFeeDetail,
		mapp.Cdc)
	keeper = NewKeeper(mapp.Cdc, keyGov, pk, tk, mapp.FeeCollectionKeeper, pk.Subspace(DefaultParamspace), ck, sk, DefaultCodespace)
	mapp.ParamsKeeper.RegisterParamSet(DefaultParamspace, &GovParams{})
	mapp.Router().AddRoute(RouterKey, NewHandler(keeper))
	mapp.QueryRouter().AddRoute(QuerierRoute, NewQuerier(keeper))

	mapp.SetEndBlocker(getEndBlocker(keeper))
	mapp.SetInitChainer(getInitChainer(mapp, keeper, sk, genState))

	require.NoError(t, mapp.CompleteSetup(keyStaking, tkeyStaking, keyGov, keyToken, keyFreeze, keyLock, keyTokenPair, keyFeeDetail))

	valTokens := sdk.NewIntFromBigInt(sdk.NewDec(420000).Int)
	if genAccs == nil || len(genAccs) == 0 {
		genAccs, addrs, pubKeys, privKeys = mock.CreateGenAccounts(numGenAccs,
			sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, valTokens)})
	}

	//mapp.MountStores(
	//	keyToken,
	//	keyFreeze,
	//	keyLock,
	//	keyTokenPair,
	//	keyFeeDetail,
	//)
	mock.SetGenesis(mapp, genAccs)

	return mapp, keeper, sk, addrs, pubKeys, privKeys
}

// gov and staking endblocker
func getEndBlocker(keeper Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		tags := EndBlocker(ctx, keeper)
		return abci.ResponseEndBlock{
			Tags: tags,
		}
	}
}

// gov and staking initchainer
func getInitChainer(mapp *mock.App, keeper Keeper, stakingKeeper staking.Keeper, genState GenesisState) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)

		stakingGenesis := staking.DefaultGenesisState()
		tokens := sdk.TokensFromTendermintPower(100000)
		stakingGenesis.Pool.NotBondedTokens = tokens

		validators, err := staking.InitGenesis(ctx, stakingKeeper, stakingGenesis)
		if err != nil {
			panic(err)
		}
		if genState.IsEmpty() {
			InitGenesis(ctx, keeper, DefaultGenesisState())
		} else {
			InitGenesis(ctx, keeper, genState)
		}
		return abci.ResponseInitChain{
			Validators: validators,
		}
	}
}

// TODO: Remove once address interface has been implemented (ref: #2186)
func SortValAddresses(addrs []sdk.ValAddress) {
	var byteAddrs [][]byte
	for _, addr := range addrs {
		byteAddrs = append(byteAddrs, addr.Bytes())
	}

	SortByteArrays(byteAddrs)

	for i, byteAddr := range byteAddrs {
		addrs[i] = byteAddr
	}
}

// Sorts Addresses
func SortAddresses(addrs []sdk.AccAddress) {
	var byteAddrs [][]byte
	for _, addr := range addrs {
		byteAddrs = append(byteAddrs, addr.Bytes())
	}
	SortByteArrays(byteAddrs)
	for i, byteAddr := range byteAddrs {
		addrs[i] = byteAddr
	}
}

// implement `Interface` in sort package.
type sortByteArrays [][]byte

func (b sortByteArrays) Len() int {
	return len(b)
}

func (b sortByteArrays) Less(i, j int) bool {
	// bytes package already implements Comparable for []byte.
	switch bytes.Compare(b[i], b[j]) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		log.Panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
		return false
	}
}

func (b sortByteArrays) Swap(i, j int) {
	b[j], b[i] = b[i], b[j]
}

// Public
func SortByteArrays(src [][]byte) [][]byte {
	sorted := sortByteArrays(src)
	sort.Sort(sorted)
	return sorted
}

func CreateTestInput(t *testing.T, isCheckTx bool) (sdk.Context, Keeper, *sdk.KVStoreKey, []byte) {
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyGov := sdk.NewKVStoreKey(StoreKey)
	keyToken := sdk.NewKVStoreKey("token")
	keyFreeze := sdk.NewKVStoreKey("freeze")
	keyLock := sdk.NewKVStoreKey("lock")
	keyTokenPair := sdk.NewKVStoreKey("token_pair")
	keyFeeDetail := sdk.NewKVStoreKey("fee_detail")

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyGov, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyToken, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyFreeze, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyLock, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyTokenPair, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyFeeDetail, sdk.StoreTypeIAVL, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, nil)

	cdc := codec.New()
	RegisterCodec(cdc)

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)

	accountKeeper := auth.NewAccountKeeper(
		cdc,    // amino codec
		keyAcc, // target store
		pk.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount, // prototype
	)

	ck := bank.NewBaseKeeper(
		accountKeeper,
		pk.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)

	sk := staking.NewKeeper(cdc, keyStaking, tkeyStaking, ck, pk.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)

	tk := token.NewKeeper(ck,
		pk,
		pk.Subspace(token.DefaultParamspace),
		auth.FeeCollectionKeeper{},
		keyToken,
		keyFreeze,
		keyLock,
		keyTokenPair,
		keyFeeDetail,
		cdc)
	keeper := NewKeeper(cdc, keyGov, pk, tk, auth.FeeCollectionKeeper{}, pk.Subspace("testgov"), ck, sk, DefaultCodespace)

	return ctx, keeper, keyParams, []byte("testgov")
}

type Codec struct {
	*codec.Codec
}

func (cdc *Codec) MarshalJSON(o interface{}) ([]byte, error) {
	return nil, fmt.Errorf("test")
}

func (cdc *Codec) UnmarshalJSON(bz []byte, ptr interface{}) error {
	return fmt.Errorf("test")
}