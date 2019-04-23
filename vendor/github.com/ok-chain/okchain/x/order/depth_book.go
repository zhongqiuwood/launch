package order

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ok-chain/okchain/x/common"
	"strings"
)

type Deal struct {
	OrderId  string  `json:"orderId"`
	Side     string  `json:"side"`
	Quantity sdk.Dec `json:"quantity"`
}

type MatchResult struct {
	Price    sdk.Dec `json:"price"`
	Quantity sdk.Dec `json:"quantity"`
	Deals    []Deal  `json:"deals"`
}

type BlockMatchResult struct {
	BlockHeight int64                  `json:"blockHeight"`
	ResultMap   map[string]MatchResult `json:"resultMap"`
	TimeStamp   int64                  `json:"timestamp"`
}

type DepthBookItem struct {
	Price        sdk.Dec `json:"price"`
	BuyQuantity  sdk.Dec `json:"buyQuantity"`
	SellQuantity sdk.Dec `json:"sellQuantity"`
}

type DepthBook []DepthBookItem

// items in depth book are sorted by price desc
// insert a new order into depth book
func (depthBook *DepthBook) InsertOrder(order *Order) {
	bookLength := len(*depthBook)
	newItem := DepthBookItem{Price: order.Price, BuyQuantity: sdk.ZeroDec(), SellQuantity: sdk.ZeroDec()}
	if order.Side == "BUY" {
		newItem.BuyQuantity = order.Quantity
	} else {
		newItem.SellQuantity = order.Quantity
	}
	if bookLength == 0 || order.Price.LT((*depthBook)[bookLength-1].Price) {
		*depthBook = append(*depthBook, newItem)
		return
	}

	index := 0
	for i, item := range *depthBook {
		if order.Price.GTE(item.Price) {
			index = i
			break
		}
	}

	if order.Price.Equal((*depthBook)[index].Price) {
		if order.Side == "BUY" {
			(*depthBook)[index].BuyQuantity = (*depthBook)[index].BuyQuantity.Add(order.RemainQuantity)
		} else {
			(*depthBook)[index].SellQuantity = (*depthBook)[index].SellQuantity.Add(order.RemainQuantity)
		}
	} else { // order.InitPrice > depthBook[index].InitPrice
		rear := append([]DepthBookItem{newItem}, (*depthBook)[index:]...)
		*depthBook = append((*depthBook)[:index], rear...)
	}
}

// remove an order from depth book when order cancelled/expired
func (depthBook *DepthBook) RemoveOrder(order *Order) {
	bookLen := len(*depthBook)
	for i := 0; i < bookLen; i++ {
		if (*depthBook)[i].Price.Equal(order.Price) {
			if order.Side == "BUY" {
				(*depthBook)[i].BuyQuantity = (*depthBook)[i].BuyQuantity.Sub(order.RemainQuantity)
			} else {
				(*depthBook)[i].SellQuantity = (*depthBook)[i].SellQuantity.Sub(order.RemainQuantity)
			}
			if (*depthBook)[i].BuyQuantity.IsZero() && (*depthBook)[i].SellQuantity.IsZero() {
				*depthBook = append((*depthBook)[:i], (*depthBook)[i+1:]...)
			}
			break
		}
	}
}

// key: product:price:side, value: orderIds
type OrderIdsMap map[string]*[]string

// insert a new order into orderIdsMap
func (orderIdsMap OrderIdsMap) InsertOrder(order *Order, ctx sdk.Context, keeper Keeper) {
	key := FormatProductPriceSideKey(order.Product, order.Price, order.Side)
	orderIds, ok := orderIdsMap[key]
	// if key not found in orderIdsMap, try get orderIds from keeper
	if !ok {
		orderIds = keeper.GetProductPriceOrderIds(ctx, key)
	}
	*orderIds = append(*orderIds, order.OrderId)
	orderIdsMap[key] = orderIds
}

// remove an order from orderIdsMap when order cancelled/expired
func (orderIdsMap OrderIdsMap) RemoveOrder(order *Order, ctx sdk.Context, keeper Keeper) {
	key := FormatProductPriceSideKey(order.Product, order.Price, order.Side)
	orderIds, ok := orderIdsMap[key]
	// if key not found in orderIdsMap, try get orderIds from keeper
	if !ok {
		orderIds = keeper.GetProductPriceOrderIds(ctx, key)
	}
	orderIdsLen := len(*orderIds)
	for i := 0; i < orderIdsLen; i++ {
		if (*orderIds)[i] == order.OrderId {
			*orderIds = append((*orderIds)[:i], (*orderIds)[i+1:]...)
			orderIdsMap[key] = orderIds
			break
		}
	}
}

// Fill orders in orderIdsMap at specific key
func (orderIdsMap OrderIdsMap) Fill(ctx sdk.Context, keeper Keeper, key string, needFillAmt sdk.Dec, fillPrice sdk.Dec, feeParams Params) *[]Deal {
	orderIds, ok := orderIdsMap[key]
	// if key not found in orderIdsMap, try get orderIds from keeper
	if !ok {
		orderIds = keeper.GetProductPriceOrderIds(ctx, key)
	}

	deals := []Deal{}
	filledAmt := sdk.ZeroDec()
	index := 0
	for filledAmt.LT(needFillAmt) {
		order := keeper.GetOrder(ctx, (*orderIds)[index])
		if filledAmt.Add(order.RemainQuantity).LTE(needFillAmt) {
			filledAmt = filledAmt.Add(order.RemainQuantity)
			deals = append(deals, Deal{order.OrderId, order.Side, order.RemainQuantity})
			fillOrder(order, ctx, keeper, fillPrice, order.RemainQuantity, feeParams)
			index++
		} else {
			deals = append(deals, Deal{order.OrderId, order.Side, needFillAmt.Sub(filledAmt)})
			fillOrder(order, ctx, keeper, fillPrice, needFillAmt.Sub(filledAmt), feeParams)
			break
		}
	}
	*orderIds = (*orderIds)[index:] // update orderIds, remove filled orderIds
	// Note: orderIds cannot be nil, we will use empty slice to remove data on keeper
	if len(*orderIds) == 0 {
		orderIdsMap[key] = &[]string{}
	}
	return &deals
}

func (orderIdsMap OrderIdsMap) SaveToKeeper(ctx sdk.Context, keeper Keeper) {
	for key, orderIds := range orderIdsMap {
		keeper.SetProductPriceOrderIds(ctx, key, orderIds)
	}
}

// Fill an order. Update order, charge fee and transfer tokens.
func fillOrder(order *Order, ctx sdk.Context, keeper Keeper, fillPrice, fillQuantity sdk.Dec, feeParams Params) {
	// update order
	order.Fill(fillPrice, fillQuantity)
	symbols := strings.Split(order.Product, "_")

	// transfer tokens
	if order.Side == "BUY" {
		burnCoins := sdk.DecCoins{{symbols[1], fillPrice.Mul(fillQuantity)}}
		receiveCoins := sdk.DecCoins{{symbols[0], fillQuantity}}
		keeper.BurnLockedCoins(ctx, order.Sender, burnCoins)
		keeper.ReceiveLockedCoins(ctx, order.Sender, receiveCoins)
	} else {
		burnCoins := sdk.DecCoins{{symbols[0], fillQuantity}}
		receiveCoins := sdk.DecCoins{{symbols[1], fillPrice.Mul(fillQuantity)}}
		keeper.BurnLockedCoins(ctx, order.Sender, burnCoins)
		keeper.ReceiveLockedCoins(ctx, order.Sender, receiveCoins)
	}

	// charge fee
	fee := GetDealFee(order, fillQuantity, ctx, keeper, feeParams)
	keeper.AddCollectedFees(ctx, fee, order.Sender, common.FeeTypeOrderDeal)
	order.RecordOrderDealFee(fee)
	keeper.SubtractCoins(ctx, order.Sender, fee)

	// update order to keeper
	keeper.SetOrder(ctx, order.OrderId, order)
	// record updated orderId
	keeper.AddUpdatedOrderId(ctx, ctx.BlockHeight(), order.OrderId)
}

// FillDepthBook will fill orders in depth book with bestPrice.
// It will update book and orderIdsMap, also update orders, charge fees, and transfer tokens, then return all deals.
func FillDepthBook(book *DepthBook, orderIdsMap OrderIdsMap, ctx sdk.Context, keeper Keeper, product string, bestPrice, maxExecution sdk.Dec, feeParams Params) *[]Deal {
	var deals []Deal
	if maxExecution.IsZero() {
		return &deals
	}
	buyAmount := sdk.ZeroDec()
	sellAmount := sdk.ZeroDec()

	// Fill buy orders, prices from high to low
	index := 0
	for index < len(*book) {
		if (*book)[index].Price.GTE(bestPrice) && buyAmount.LT(maxExecution) { // item.InitPrice >= bestPrice, fill buy orders
			fillAmount := sdk.MinDec((*book)[index].BuyQuantity, maxExecution.Sub(buyAmount))
			if fillAmount.IsZero() {
				index++
				continue
			}
			buyAmount = buyAmount.Add(fillAmount)
			(*book)[index].BuyQuantity = (*book)[index].BuyQuantity.Sub(fillAmount)

			// Fill buy orders at this price
			key := FormatProductPriceSideKey(product, (*book)[index].Price, "BUY")
			deals = append(deals, *(orderIdsMap.Fill(ctx, keeper, key, fillAmount, bestPrice, feeParams))...)

			// If this item has no buy quantity or sell quantity anymore, remove it from depth book.
			if (*book)[index].BuyQuantity.IsZero() && (*book)[index].SellQuantity.IsZero() {
				*book = append((*book)[:index], (*book)[index+1:]...)
			} else {
				index++
			}
		} else {
			break
		}
	}
	// Fill sell orders, prices from low to high
	index = len(*book) - 1
	for index >= 0 {
		if (*book)[index].Price.LTE(bestPrice) && sellAmount.LT(maxExecution) { // item.InitPrice <= bestPrice, fill sell orders
			fillAmount := sdk.MinDec((*book)[index].SellQuantity, maxExecution.Sub(sellAmount))
			if fillAmount.IsZero() {
				index--
				continue
			}
			sellAmount = sellAmount.Add(fillAmount)
			(*book)[index].SellQuantity = (*book)[index].SellQuantity.Sub(fillAmount)

			// Fill sell orders at this price
			key := FormatProductPriceSideKey(product, (*book)[index].Price, "SELL")
			deals = append(deals, *(orderIdsMap.Fill(ctx, keeper, key, fillAmount, bestPrice, feeParams))...)

			// If this item has no buy quantity or sell quantity anymore, remove it from depth book.
			if (*book)[index].BuyQuantity.IsZero() && (*book)[index].SellQuantity.IsZero() {
				*book = append((*book)[:index], (*book)[index+1:]...)
			}
			index--
		} else {
			break
		}
	}
	return &deals
}
