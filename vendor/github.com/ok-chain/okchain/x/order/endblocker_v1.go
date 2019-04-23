package order

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ok-chain/okchain/util"
	"github.com/ok-chain/okchain/x/common"
	"github.com/ok-chain/okchain/x/version"

	//"github.com/ok-chain/okchain/x/version"
)

const (
	OrderExpireBlocks = 86400     // if 1 block per second, order will be expired after 24 hours
	DataExpireBlocks  = 86400 * 3 // if 1 block per second, some data will be expired after 3 days
)

// Called every block, check expired orders
func EndBlocker(ctx sdk.Context, keeper Keeper) sdk.Tags {
	return version.GetVersionController().EndBlocker(ctx, keeper, getModule())
	//resTags := sdk.NewTags()
	//return resTags
}

// Called every block, check expired orders
func endBlockerV1(ctx sdk.Context, keeper Keeper) sdk.Tags {
	expireTags := checkExpiredOrders(ctx, keeper)
	matchTags := allPeriodicAuctionMatch(ctx, keeper)
	dropExpiredData(ctx, keeper)
	resTags := expireTags.AppendTags(matchTags)
	return resTags
}


func checkExpiredOrders(ctx sdk.Context, keeper Keeper) sdk.Tags {
	logger := ctx.Logger().With("module", "x/order")
	resTags := sdk.NewTags()
	blockHeight := blockHeight(ctx)
	expiredBlockHeight := blockHeight - OrderExpireBlocks
	feeParams := keeper.GetParams(ctx)
	orderIdsMap := make(OrderIdsMap)

	// check orders in expired block, remove expired orders by order id
	orderNum := keeper.GetBlockOrderNum(ctx, expiredBlockHeight)
	var index int64 = 1
	for ; index < orderNum; index++ {
		orderId := common.FormatOrderId(expiredBlockHeight, index)

		order := keeper.GetOrder(ctx, orderId)
		if order != nil && order.Status == OrderStatusOpen {
			// update order
			order.Expire()
			logger.Info(fmt.Sprintf("order (%s) expired", order.OrderId))
			// unlock coins in this order & charge fee
			needUnlockCoins := order.NeedUnlockCoins()
			keeper.UnlockCoins(ctx, order.Sender, needUnlockCoins)
			if order.Status == OrderStatusExpired { // charge fees only if fully expired
				fee := GetOrderExpireFee(order, ctx, keeper, feeParams)
				keeper.AddCollectedFees(ctx, fee, order.Sender, common.FeeTypeOrderExpire)
				order.RecordOrderExpireFee(fee)
				keeper.SubtractCoins(ctx, order.Sender, fee)
			}
			keeper.SetOrder(ctx, orderId, order)

			// update depth book and orderIdsMap
			depthBook := keeper.GetDepthBook(ctx, order.Product)
			depthBook.RemoveOrder(order)
			orderIdsMap.RemoveOrder(order, ctx, keeper)
			keeper.SetDepthBook(ctx, order.Product, depthBook)
			keeper.AddUpdatedOrderId(ctx, blockHeight, order.OrderId)
		}
	}
	orderIdsMap.SaveToKeeper(ctx, keeper)
	return resTags
}

func dump(ctx sdk.Context, format string, a ...interface{}) {
	logger := ctx.Logger().With("module", "x/order")
	format = fmt.Sprintf("[%s]%s", util.GoId, format)
	logger.Info(fmt.Sprintf(format, a...))
}

func allPeriodicAuctionMatch(ctx sdk.Context, keeper Keeper) sdk.Tags {
	//logger := ctx.Logger().With("module", "x/order")
	resTags := sdk.NewTags()
	blockHeight := blockHeight(ctx)
	orderNum := keeper.GetBlockOrderNum(ctx, blockHeight)
	depthBookMap := make(map[string]*DepthBook)
	orderIdsMap := make(OrderIdsMap)
	var i int64
	// step0: handle new orders in current block, insert into depth book and orderIdsMap
	for i = 1; i < orderNum; i++ {
		orderId := common.FormatOrderId(blockHeight, i)
		order := keeper.GetOrder(ctx, orderId)
		if depthBook, ok := depthBookMap[order.Product]; ok {
			depthBook.InsertOrder(order)
		} else {
			newDepthBook := keeper.GetDepthBook(ctx, order.Product)
			newDepthBook.InsertOrder(order)
			depthBookMap[order.Product] = newDepthBook
		}
		orderIdsMap.InsertOrder(order, ctx, keeper)
	}

	// step1: execute periodic auction match, get best price and max execution quantity, save latest price
	resultMap := make(map[string]MatchResult)
	for product, depthBook := range depthBookMap {
		tokenPair := keeper.tokenKeeper.GetTokenPair(ctx, product)
		bestPrice, maxExecution := periodicAuctionMatchPrice(depthBook, int64(tokenPair.MaxPriceDigit), keeper.GetLastPrice(ctx, product))
		if maxExecution.IsPositive() {
			keeper.SetLastPrice(ctx, product, bestPrice)
			resultMap[product] = MatchResult{Price: bestPrice, Quantity: maxExecution, Deals: []Deal{}}
		}
	}

	// step2: fill orders in match results, transfer tokens and collect fees
	feeParams := keeper.GetParams(ctx)
	for product, matchResult := range resultMap {
		deals := FillDepthBook(depthBookMap[product], orderIdsMap, ctx, keeper, product, matchResult.Price, matchResult.Quantity, feeParams)
		matchResult.Deals = *deals
		resultMap[product] = matchResult
	}

	// step3: save new depthBook and orderIdsMap
	for product, depthBook := range depthBookMap {
		keeper.SetDepthBook(ctx, product, depthBook)
	}
	orderIdsMap.SaveToKeeper(ctx, keeper)

	// step4: save match results for querying
	if len(resultMap) > 0 {
		blockMatchResult := &BlockMatchResult{
			blockHeight,
			resultMap,
			ctx.BlockHeader().Time.Unix(),
		}
		keeper.SetBlockMatchResult(ctx, blockHeight, blockMatchResult)
	}
	return resTags
}

func dropExpiredData(ctx sdk.Context, keeper Keeper) {
	blockHeight := blockHeight(ctx)
	expiredBlockHeight := blockHeight - DataExpireBlocks

	// drop orders at expiredBlock
	orderNum := keeper.GetBlockOrderNum(ctx, expiredBlockHeight)
	var index int64 = 1
	for ; index < orderNum; index++ {
		orderId := common.FormatOrderId(expiredBlockHeight, index)
		keeper.DropOrder(ctx, orderId)
	}
	// drop updated order ids at expiredBlock
	keeper.DropUpdatedOrderIds(ctx, expiredBlockHeight)
	// drop block order num at expiredBlock
	keeper.DropBlockOrderNum(ctx, expiredBlockHeight)
	// drop block match result at expiredBlock
	keeper.DropBlockMatchResult(ctx, expiredBlockHeight)
}
