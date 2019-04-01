package order

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"log"
	"strings"
)

const (
	OrderExpireBlocks = 86400 // if 1 block per second, order will be expired after 24 hours
)

// Called every block, check expired orders
func EndBlocker(ctx sdk.Context, keeper Keeper) sdk.Tags {
	expireTags := checkExpiredOrders(ctx, keeper)
	matchTags := allPeriodicAuctionMatch(ctx, keeper)
	resTags := expireTags.AppendTags(matchTags)
	return resTags
}

func checkExpiredOrders(ctx sdk.Context, keeper Keeper) sdk.Tags {
	logger := ctx.Logger().With("module", "x/order")
	resTags := sdk.NewTags()
	blockHeight := ctx.BlockHeight()
	expiredBlockHeight := blockHeight - OrderExpireBlocks

	// remove expired orders by order id
	orderNum := keeper.GetBlockOrderNum(ctx, expiredBlockHeight)
	var index int64 = 0
	for ; index < orderNum; index++ {
		orderId := FormatOrderId(blockHeight, index)
		order := keeper.GetOrder(ctx, orderId)
		if order != nil && order.Status == OrderStatusOpen {
			// update order
			order.Expire()
			keeper.SetOrder(ctx, orderId, order)
			logger.Info(fmt.Sprintf("order (%s) expired", order.OrderId))
			// unlock coins in this order & charge fee
			needUnlockCoins := order.NeedUnlockCoins()
			fee, inOrder := GetOrderExpireFee(order, ctx, keeper)
			keeper.AddCollectedFees(ctx, fee)
			if inOrder {
				needUnlockCoins = needUnlockCoins.Sub(fee)
				keeper.BurnLockedCoins(ctx, order.Sender, fee)
			} else {
				keeper.SubtractCoins(ctx, order.Sender, fee)
			}
			keeper.UnlockCoins(ctx, order.Sender, needUnlockCoins)

			// update depth book
			depthBook := keeper.GetDepthBook(ctx, order.Product)
			depthBook.RemoveOrder(order)
			keeper.SetDepthBook(ctx, order.Product, depthBook)
		}
	}
	return resTags
}

func allPeriodicAuctionMatch(ctx sdk.Context, keeper Keeper) sdk.Tags {
	logger := ctx.Logger().With("module", "x/order")
	resTags := sdk.NewTags()
	blockHeight := ctx.BlockHeight()
	orderNum := keeper.GetBlockOrderNum(ctx, blockHeight)
	depthBookMap := make(map[string]*DepthBook)
	var i int64
	for i = 0; i < orderNum; i++ {
		orderId := FormatOrderId(blockHeight, i)
		order := keeper.GetOrder(ctx, orderId)
		if depthBook, ok := depthBookMap[order.Product]; ok {
			depthBook.InsertOrder(order)
		} else {
			newDepthBook := keeper.GetDepthBook(ctx, order.Product)
			newDepthBook.InsertOrder(order)
			depthBookMap[order.Product] = newDepthBook
		}
	}

	// execute periodic auction match
	resultMap := make(map[string]MatchResult)
	for product, depthBook := range depthBookMap {
		// TODO: use the precision of token pair
		matchResult := PeriodicAuctionMatch(depthBook, 1, keeper.GetLastPrice(ctx, product))
		keeper.SetDepthBook(ctx, product, depthBook)
		if matchResult.Quantity.IsPositive() {
			keeper.SetLastPrice(ctx, product, matchResult.Price)
			resultMap[product] = *matchResult
		}
		logger.Info(fmt.Sprintf("filled match result (%s-%s): %v", blockHeight, product, matchResult))
		log.Printf("filled match result (%d-%s): %v\n", blockHeight, product, matchResult)
	}

	// fill orders, transfer tokens and collect fees
	for _, matchResult := range resultMap {
		// update orders
		for _, deal := range matchResult.Deals {
			order := keeper.GetOrder(ctx, deal.OrderId)
			order.Fill(deal.Quantity)
			symbols := strings.Split(order.Product, "_")
			// charge fee
			fee, inOrder := GetDealFee(&deal, order, ctx, keeper)
			keeper.AddCollectedFees(ctx, fee)
			if !inOrder {
				keeper.SubtractCoins(ctx, order.Sender, fee)
			}
			// transfer tokens
			if order.Side == "BUY" {
				burnCoins := sdk.DecCoins{{symbols[1], matchResult.Price.Mul(deal.Quantity)}}
				receiveCoins := sdk.DecCoins{{symbols[0], deal.Quantity}}
				if inOrder {
					receiveCoins = receiveCoins.Sub(fee)
				}
				keeper.BurnLockedCoins(ctx, order.Sender, burnCoins)
				keeper.ReceiveLockedCoins(ctx, order.Sender, receiveCoins)
			} else {
				burnCoins := sdk.DecCoins{{symbols[0], deal.Quantity}}
				receiveCoins := sdk.DecCoins{{symbols[1], matchResult.Price.Mul(deal.Quantity)}}
				if inOrder {
					receiveCoins = receiveCoins.Sub(fee)
				}
				keeper.BurnLockedCoins(ctx, order.Sender, burnCoins)
				keeper.ReceiveLockedCoins(ctx, order.Sender, receiveCoins)
			}
			keeper.SetOrder(ctx, order.OrderId, order)
		}
	}

	if len(resultMap) > 0 {
		keeper.SetMatchResultMap(ctx, blockHeight, resultMap)
	}
	return resTags
}
