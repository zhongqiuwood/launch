package order

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	// The reference to the TokenKeeper to modify balances
	tokenKeeper TokenKeeper
	// The reference to the Param Keeper to get and set Global Params
	paramsKeeper params.Keeper
	// The reference to the Paramstore to get and set gov specific params
	paramSpace params.Subspace
	// The reference to the FeeCollectionKeeper to collect fees
	feeCollectionKeeper auth.FeeCollectionKeeper

	ordersStoreKey sdk.StoreKey // Unexposed key to access name store from sdk.Context
	tradesStoreKey sdk.StoreKey
	cdc            *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(tokenKeeper TokenKeeper, paramsKeeper params.Keeper, paramSpace params.Subspace, feeCollectionKeeper auth.FeeCollectionKeeper, ordersStoreKey, tradesStoreKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		tokenKeeper:         tokenKeeper,
		paramsKeeper:        paramsKeeper,
		paramSpace:          paramSpace.WithKeyTable(ParamKeyTable()),
		feeCollectionKeeper: feeCollectionKeeper,
		ordersStoreKey:      ordersStoreKey,
		tradesStoreKey:      tradesStoreKey,
		cdc:                 cdc,
	}
}

// NewOrder - place an order
func (k Keeper) SetOrder(ctx sdk.Context, orderId string, order *Order) {
	store := ctx.KVStore(k.ordersStoreKey)
	store.Set([]byte(orderId), k.cdc.MustMarshalBinaryBare(order))
}

func (k Keeper) GetOrder(ctx sdk.Context, orderId string) *Order {
	store := ctx.KVStore(k.ordersStoreKey)
	orderInfo := store.Get([]byte(orderId))
	if orderInfo == nil {
		return nil
	}
	order := &Order{}
	k.cdc.MustUnmarshalBinaryBare(orderInfo, order)
	return order
}

func (k Keeper) SetDepthBook(ctx sdk.Context, product string, depthBook *DepthBook) {
	store := ctx.KVStore(k.ordersStoreKey)
	depthBookKey := fmt.Sprintf("depthbook:%v", product)
	if depthBook == nil || len(*depthBook) == 0 {
		store.Delete([]byte(depthBookKey))
	} else {
		store.Set([]byte(depthBookKey), k.cdc.MustMarshalBinaryBare(depthBook))
	}
}

func (k Keeper) GetDepthBook(ctx sdk.Context, product string) *DepthBook {
	store := ctx.KVStore(k.ordersStoreKey)
	depthBookKey := fmt.Sprintf("depthbook:%v", product)
	bookBytes := store.Get([]byte(depthBookKey))
	if bookBytes == nil {
		// Return an empty DepthBook instead of nil
		return &DepthBook{}
	}
	depthBook := &DepthBook{}
	k.cdc.MustUnmarshalBinaryBare(bookBytes, depthBook)
	return depthBook
}

func (k Keeper) SetLastPrice(ctx sdk.Context, product string, price sdk.Dec) {
	store := ctx.KVStore(k.ordersStoreKey)
	lastPriceKey := fmt.Sprintf("lastprice:%v", product)
	store.Set([]byte(lastPriceKey), k.cdc.MustMarshalBinaryBare(price))
}

func (k Keeper) GetLastPrice(ctx sdk.Context, product string) sdk.Dec {
	store := ctx.KVStore(k.ordersStoreKey)
	lastPriceKey := fmt.Sprintf("lastprice:%v", product)
	priceBytes := store.Get([]byte(lastPriceKey))
	if priceBytes == nil {
		// Return an empty OrderBook instead of nil
		return sdk.ZeroDec()
	}
	var price sdk.Dec
	k.cdc.MustUnmarshalBinaryBare(priceBytes, &price)
	return price
}

// get the num of orders in specific block
func (k Keeper) GetBlockOrderNum(ctx sdk.Context, blockHeight int64) int64 {
	store := ctx.KVStore(k.ordersStoreKey)
	key := fmt.Sprintf("orderNum:block(%v)", blockHeight)
	numBytes := store.Get([]byte(key))
	if numBytes == nil {
		return 0
	}
	return BytesToInt64(numBytes)
}

func (k Keeper) SetBlockOrderNum(ctx sdk.Context, blockHeight int64, orderNum int64) {
	store := ctx.KVStore(k.ordersStoreKey)
	key := fmt.Sprintf("orderNum:block(%v)", blockHeight)
	store.Set([]byte(key), Int64ToBytes(orderNum))
}

// use tradesStoreKey
func (k Keeper) SetMatchResultMap(ctx sdk.Context, blockHeight int64, matchResultMap map[string]MatchResult) {
	store := ctx.KVStore(k.tradesStoreKey)
	key := fmt.Sprintf("matchResultMap:%d", blockHeight)
	if matchResultMap == nil || len(matchResultMap) == 0 {
		store.Delete([]byte(key))
	} else {
		store.Set([]byte(key), k.cdc.MustMarshalJSON(matchResultMap))
	}
}

func (k Keeper) GetMatchResultMap(ctx sdk.Context, blockHeight int64) map[string]MatchResult {
	store := ctx.KVStore(k.tradesStoreKey)
	key := fmt.Sprintf("matchResultMap:%d", blockHeight)
	matchResultBytes := store.Get([]byte(key))
	matchResultMap := make(map[string]MatchResult)
	if matchResultBytes == nil {
		return matchResultMap
	}
	k.cdc.MustUnmarshalJSON(matchResultBytes, &matchResultMap)
	return matchResultMap
}

// use tokenKeeper
func (k Keeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins) bool {
	baseCoins := ConvertDecCoinsToCoins(coins)
	return k.tokenKeeper.HasCoins(ctx, addr, baseCoins)
}

func (k Keeper) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins) error {
	baseCoins := ConvertDecCoinsToCoins(coins)
	return k.tokenKeeper.SubtractCoins(ctx, addr, baseCoins)
}

func (k Keeper) LockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins) error {
	baseCoins := ConvertDecCoinsToCoins(coins)
	return k.tokenKeeper.LockCoins(ctx, addr, baseCoins)
}

func (k Keeper) UnlockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins) error {
	baseCoins := ConvertDecCoinsToCoins(coins)
	return k.tokenKeeper.UnlockCoins(ctx, addr, baseCoins)
}

func (k Keeper) BurnLockedCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins) error {
	baseCoins := ConvertDecCoinsToCoins(coins)
	return k.tokenKeeper.BurnLockedCoins(ctx, addr, baseCoins)
}

func (k Keeper) ReceiveLockedCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins) error {
	baseCoins := ConvertDecCoinsToCoins(coins)
	return k.tokenKeeper.ReceiveLockedCoins(ctx, addr, baseCoins)
}

// use feeCollectionKeeper
// AddCollectedFees - add to the fee pool
func (k Keeper) AddCollectedFees(ctx sdk.Context, coins sdk.DecCoins) sdk.Coins {
	baseCoins := ConvertDecCoinsToCoins(coins)
	return k.feeCollectionKeeper.AddCollectedFees(ctx, baseCoins)
}

// get inflation params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) (params Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// set inflation params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
