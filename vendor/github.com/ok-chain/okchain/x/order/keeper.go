package order

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/ok-chain/okchain/x/common"
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

func (k Keeper) DropOrder(ctx sdk.Context, orderId string) {
	store := ctx.KVStore(k.ordersStoreKey)
	store.Delete([]byte(orderId))
}

func (k Keeper) GetUpdatedOrderIds(ctx sdk.Context, blockHeight int64) []string {
	store := ctx.KVStore(k.ordersStoreKey)
	key := fmt.Sprintf("updatedAt(%v)", blockHeight)
	bytes := store.Get([]byte(key))
	var orderIds []string
	if bytes == nil {
		return orderIds
	}
	k.cdc.MustUnmarshalJSON(bytes, &orderIds)
	return orderIds
}

func (k Keeper) AddUpdatedOrderId(ctx sdk.Context, blockHeight int64, orderId string) {
	store := ctx.KVStore(k.ordersStoreKey)
	key := fmt.Sprintf("updatedAt(%v)", blockHeight)
	orderIds := k.GetUpdatedOrderIds(ctx, blockHeight)
	orderIds = append(orderIds, orderId)
	store.Set([]byte(key), k.cdc.MustMarshalJSON(orderIds))
}

func (k Keeper) DropUpdatedOrderIds(ctx sdk.Context, blockHeight int64) {
	store := ctx.KVStore(k.ordersStoreKey)
	key := fmt.Sprintf("updatedAt(%v)", blockHeight)
	store.Delete([]byte(key))
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

func FormatProductPriceSideKey(product string, price sdk.Dec, side string) string {
	return fmt.Sprintf("%v:%v:%v", product, price.String(), side)
}

func (k Keeper) SetProductPriceOrderIds(ctx sdk.Context, key string, orderIds *[]string) {
	store := ctx.KVStore(k.ordersStoreKey)
	if orderIds == nil || len(*orderIds) == 0 {
		store.Delete([]byte(key))
	} else {
		store.Set([]byte(key), k.cdc.MustMarshalJSON(*orderIds))
	}
}

func (k Keeper) GetProductPriceOrderIds(ctx sdk.Context, key string) *[]string {
	store := ctx.KVStore(k.ordersStoreKey)
	bz := store.Get([]byte(key))
	orderIds := []string{}
	if bz == nil {
		return &orderIds
	}
	k.cdc.MustUnmarshalJSON(bz, &orderIds)
	return &orderIds
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
		// If last price does not exist, set the init price of token pair as last price
		tokenPair := k.tokenKeeper.GetTokenPair(ctx, product)
		k.SetLastPrice(ctx, product, tokenPair.InitPrice)
		return tokenPair.InitPrice
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
		return 1
	}
	return common.BytesToInt64(numBytes)
}

func (k Keeper) SetBlockOrderNum(ctx sdk.Context, blockHeight int64, orderNum int64) {
	store := ctx.KVStore(k.ordersStoreKey)
	key := fmt.Sprintf("orderNum:block(%v)", blockHeight)
	store.Set([]byte(key), common.Int64ToBytes(orderNum))
}

func (k Keeper) DropBlockOrderNum(ctx sdk.Context, blockHeight int64) {
	store := ctx.KVStore(k.ordersStoreKey)
	key := fmt.Sprintf("orderNum:block(%v)", blockHeight)
	store.Delete([]byte(key))
}

// use tradesStoreKey
func (k Keeper) SetBlockMatchResult(ctx sdk.Context, blockHeight int64, blockMatchResult *BlockMatchResult) {
	store := ctx.KVStore(k.tradesStoreKey)
	key := fmt.Sprintf("blockMatchResult:%d", blockHeight)
	store.Set([]byte(key), k.cdc.MustMarshalJSON(blockMatchResult))
}

func (k Keeper) GetBlockMatchResult(ctx sdk.Context, blockHeight int64) *BlockMatchResult {
	store := ctx.KVStore(k.tradesStoreKey)
	key := fmt.Sprintf("blockMatchResult:%d", blockHeight)
	bytes := store.Get([]byte(key))
	if bytes == nil {
		return nil
	}
	blockMatchResult := &BlockMatchResult{}
	k.cdc.MustUnmarshalJSON(bytes, &blockMatchResult)
	return blockMatchResult
}

func (k Keeper) DropBlockMatchResult(ctx sdk.Context, blockHeight int64) {
	store := ctx.KVStore(k.tradesStoreKey)
	key := fmt.Sprintf("blockMatchResult:%d", blockHeight)
	store.Delete([]byte(key))
}

// use tokenKeeper
func (k Keeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins) bool {
	baseCoins := common.ConvertDecCoinsToCoins(coins)
	return k.tokenKeeper.HasCoins(ctx, addr, baseCoins)
}

func (k Keeper) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins) error {
	baseCoins := common.ConvertDecCoinsToCoins(coins)
	return k.tokenKeeper.SubtractCoins(ctx, addr, baseCoins)
}

func (k Keeper) LockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins) error {
	baseCoins := common.ConvertDecCoinsToCoins(coins)
	return k.tokenKeeper.LockCoins(ctx, addr, baseCoins)
}

func (k Keeper) UnlockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins) error {
	baseCoins := common.ConvertDecCoinsToCoins(coins)
	return k.tokenKeeper.UnlockCoins(ctx, addr, baseCoins)
}

func (k Keeper) BurnLockedCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins) error {
	baseCoins := common.ConvertDecCoinsToCoins(coins)
	return k.tokenKeeper.BurnLockedCoins(ctx, addr, baseCoins)
}

func (k Keeper) ReceiveLockedCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins) error {
	baseCoins := common.ConvertDecCoinsToCoins(coins)
	return k.tokenKeeper.ReceiveLockedCoins(ctx, addr, baseCoins)
}

// use feeCollectionKeeper
// AddCollectedFees - add to the fee pool
func (k Keeper) AddCollectedFees(ctx sdk.Context, coins sdk.DecCoins, from sdk.AccAddress, feeType string) sdk.Coins {
	k.tokenKeeper.AddFeeDetail(ctx, from.String(), coins.String(), feeType)
	baseCoins := common.ConvertDecCoinsToCoins(coins)
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
