package order

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
)

type OrderItem struct {
	OrderId        string  `json:"orderId"`
	RemainQuantity sdk.Dec `json:"quantity"`
}

type DepthBookItem struct {
	Price        sdk.Dec     `json:"price"`
	BuyQuantity  sdk.Dec     `json:"buyQuantity"`
	SellQuantity sdk.Dec     `json:"sellQuantity"`
	BuyOrders    []OrderItem `json:"buyOrders"`
	SellOrders   []OrderItem `json:"sellOrders"`
}

type DepthBook []DepthBookItem

// items in depth book are sorted by price desc
func (depthBook *DepthBook) InsertOrder(order *Order) {
	bookLength := len(*depthBook)
	newItem := DepthBookItem{Price: order.Price, BuyQuantity: sdk.ZeroDec(), SellQuantity: sdk.ZeroDec()}
	if order.Side == "BUY" {
		newItem.BuyQuantity = order.Quantity
		newItem.BuyOrders = []OrderItem{{order.OrderId, order.RemainQuantity}}
	} else {
		newItem.SellQuantity = order.Quantity
		newItem.SellOrders = []OrderItem{{order.OrderId, order.RemainQuantity}}
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
			(*depthBook)[index].BuyOrders = append((*depthBook)[index].BuyOrders, OrderItem{order.OrderId, order.RemainQuantity})
		} else {
			(*depthBook)[index].SellQuantity = (*depthBook)[index].SellQuantity.Add(order.RemainQuantity)
			(*depthBook)[index].SellOrders = append((*depthBook)[index].SellOrders, OrderItem{order.OrderId, order.RemainQuantity})
		}
	} else { // order.Price > depthBook[index].Price
		rear := append([]DepthBookItem{newItem}, (*depthBook)[index:]...)
		*depthBook = append((*depthBook)[:index], rear...)
	}
}

func (depthBook *DepthBook) RemoveOrder(order *Order) {
	bookLen := len(*depthBook)
	for i := 0; i < bookLen; i++ {
		if (*depthBook)[i].Price.Equal(order.Price) {
			if order.Side == "BUY" {
				(*depthBook)[i].BuyQuantity = (*depthBook)[i].BuyQuantity.Sub(order.RemainQuantity)
				for j := 0; j < len((*depthBook)[i].BuyOrders); j++ {
					if (*depthBook)[i].BuyOrders[j].OrderId == order.OrderId {
						(*depthBook)[i].BuyOrders = append((*depthBook)[i].BuyOrders[:j], (*depthBook)[i].BuyOrders[j+1:]...)
						break
					}
				}
			} else {
				(*depthBook)[i].SellQuantity = (*depthBook)[i].SellQuantity.Sub(order.RemainQuantity)
				for j := 0; j < len((*depthBook)[i].SellOrders); j++ {
					if (*depthBook)[i].SellOrders[j].OrderId == order.OrderId {
						(*depthBook)[i].SellOrders = append((*depthBook)[i].SellOrders[:j], (*depthBook)[i].SellOrders[j+1:]...)
						break
					}
				}
			}
			if (*depthBook)[i].BuyQuantity.IsZero() && (*depthBook)[i].SellQuantity.IsZero() {
				*depthBook = append((*depthBook)[:i], (*depthBook)[i+1:]...)
			}
			break
		}
	}
}

// Round a decimal with precision, perform bankers rounding (gaussian rounding)
func RoundDecimal(dec sdk.Dec, precision int64) sdk.Dec {
	precisionMul := sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(precision), nil))
	return sdk.NewDecFromInt(dec.MulInt(precisionMul).RoundInt()).QuoInt(precisionMul)
}

// Calculate periodic auction match price, return the best price and execution amount
func periodicAuctionMatchPrice(book *DepthBook, pricePrecision int64, refPrice sdk.Dec) (bestPrice sdk.Dec, maxExecution sdk.Dec) {
	bookLength := len(*book)
	buyAmountSum := make([]sdk.Dec, bookLength)
	sellAmountSum := make([]sdk.Dec, bookLength)

	buyAmountSum[0] = (*book)[0].BuyQuantity
	for i := 1; i < bookLength; i++ {
		buyAmountSum[i] = buyAmountSum[i-1].Add((*book)[i].BuyQuantity)
	}

	sellAmountSum[bookLength-1] = (*book)[bookLength-1].SellQuantity
	for i := bookLength - 2; i >= 0; i-- {
		sellAmountSum[i] = sellAmountSum[i+1].Add((*book)[i].SellQuantity)
	}

	maxExecution = sdk.ZeroDec()
	execution := make([]sdk.Dec, bookLength)
	for i := 0; i < bookLength; i++ {
		execution[i] = sdk.MinDec(buyAmountSum[i], sellAmountSum[i])
		maxExecution = sdk.MaxDec(execution[i], maxExecution)
	}

	// rule 0
	if maxExecution.IsZero() {
		return refPrice, maxExecution
	}

	var indexesRule1 []int
	for i := 0; i < bookLength; i++ {
		if execution[i].Equal(maxExecution) {
			indexesRule1 = append(indexesRule1, i)
		}
	}
	indexLen1 := len(indexesRule1)
	// rule1: maximum matched volume
	if indexLen1 == 1 {
		bestPrice = (*book)[indexesRule1[0]].Price
		return
	}

	// rule2: minimum surplus
	imbalance := make([]sdk.Dec, bookLength)
	for i := 0; i < bookLength; i++ {
		imbalance[i] = buyAmountSum[i].Sub(sellAmountSum[i])
	}
	minAbsImbalance := imbalance[0].Abs()
	for i := 1; i < bookLength; i++ {
		minAbsImbalance = sdk.MinDec(minAbsImbalance, imbalance[i].Abs())
	}
	var indexesRule2 []int
	for i := 0; i < indexLen1; i++ {
		if imbalance[indexesRule1[i]].Abs().Equal(minAbsImbalance) {
			indexesRule2 = append(indexesRule2, indexesRule1[i])
		}
	}
	indexLen2 := len(indexesRule2)
	if indexLen2 == 1 {
		bestPrice = (*book)[indexesRule2[0]].Price
		return
	}

	// rule3: Market Pressure
	if imbalance[indexesRule2[0]].GT(sdk.ZeroDec()) { // rule3a: all imbalance > 0, buyer pressure
		newRefPrice := refPrice.Mul(sdk.MustNewDecFromStr("1.05"))
		newRefPrice = RoundDecimal(newRefPrice, pricePrecision)
		if (*book)[indexesRule2[0]].Price.LT(newRefPrice) { // all price < newRefPrice
			bestPrice = (*book)[indexesRule2[0]].Price
		} else if (*book)[indexesRule2[indexLen2-1]].Price.GT(newRefPrice) { // all price > newRefPrice
			bestPrice = (*book)[indexesRule2[indexLen2-1]].Price
		} else {
			bestPrice = newRefPrice
		}
	} else if imbalance[indexesRule2[indexLen2-1]].LT(sdk.ZeroDec()) { // rule3b: all imbalance < 0, seller pressure
		newRefPrice := refPrice.Mul(sdk.MustNewDecFromStr("0.95"))
		newRefPrice = RoundDecimal(newRefPrice, pricePrecision)
		if (*book)[indexesRule2[0]].Price.LT(newRefPrice) { // all price < newRefPrice
			bestPrice = (*book)[indexesRule2[0]].Price
		} else if (*book)[indexesRule2[indexLen2-1]].Price.GT(newRefPrice) { // all price > newRefPrice
			bestPrice = (*book)[indexesRule2[indexLen2-1]].Price
		} else {
			bestPrice = newRefPrice
		}
	} else { // rule3c: some imbalance > 0, and some imbalance < 0, no buyer pressure or seller pressure
		if (*book)[indexesRule2[0]].Price.LT(refPrice) { // all price < refPrice
			bestPrice = (*book)[indexesRule2[0]].Price
		} else if (*book)[indexesRule2[indexLen2-1]].Price.GT(refPrice) { // all price > refPrice
			bestPrice = (*book)[indexesRule2[indexLen2-1]].Price
		} else {
			bestPrice = RoundDecimal(refPrice, pricePrecision)
		}
	}
	return
}

type Deal struct {
	OrderId  string  `json:"orderId"`
	Side     string  `json:"side"`
	Quantity sdk.Dec `json:"quantity"`
}

type MatchResult struct {
	Price    sdk.Dec `json:"price"`
	Quantity sdk.Dec `json:"quantity"`
	Deals    []Deal  `json:"trades"`
}

func fillDepthBookItem(item *DepthBookItem, buyQuantity, sellQuantity sdk.Dec) []Deal {
	if buyQuantity.GT(item.BuyQuantity) || sellQuantity.GT(item.SellQuantity) {
		panic("fill quantity is invalid for the depth book item")
	}
	deals := []Deal{}
	filledBuyQuantity := sdk.ZeroDec()
	filledSellQuantity := sdk.ZeroDec()

	if item.BuyQuantity.Equal(buyQuantity) {
		filledBuyQuantity = buyQuantity
		for _, orderItem := range item.BuyOrders {
			deals = append(deals, Deal{orderItem.OrderId, "BUY", orderItem.RemainQuantity})
		}
		item.BuyQuantity = sdk.ZeroDec()
		item.BuyOrders = []OrderItem{}
	}
	if item.SellQuantity.Equal(sellQuantity) {
		filledSellQuantity = sellQuantity
		for _, orderItem := range item.SellOrders {
			deals = append(deals, Deal{orderItem.OrderId, "SELL", orderItem.RemainQuantity})
		}
		item.SellQuantity = sdk.ZeroDec()
		item.SellOrders = []OrderItem{}
	}

	if filledBuyQuantity.Equal(buyQuantity) && filledSellQuantity.Equal(sellQuantity) {
		return deals
	}

	if filledBuyQuantity.LT(buyQuantity) {
		item.BuyQuantity = item.BuyQuantity.Sub(buyQuantity)
		for filledBuyQuantity.LT(buyQuantity) {
			if filledBuyQuantity.Add(item.BuyOrders[0].RemainQuantity).LTE(buyQuantity) {
				filledBuyQuantity = filledBuyQuantity.Add(item.BuyOrders[0].RemainQuantity)
				deals = append(deals, Deal{item.BuyOrders[0].OrderId, "BUY", item.BuyOrders[0].RemainQuantity})
				item.BuyOrders = item.BuyOrders[1:]
			} else {
				deals = append(deals, Deal{item.BuyOrders[0].OrderId, "BUY", buyQuantity.Sub(filledBuyQuantity)})
				item.BuyOrders[0].RemainQuantity = item.BuyOrders[0].RemainQuantity.Sub(buyQuantity.Sub(filledBuyQuantity))
				break
			}
		}
	}
	if filledSellQuantity.LT(sellQuantity) {
		item.SellQuantity = item.SellQuantity.Sub(sellQuantity)
		for filledSellQuantity.LT(sellQuantity) {
			if filledSellQuantity.Add(item.SellOrders[0].RemainQuantity).LTE(sellQuantity) {
				filledSellQuantity = filledSellQuantity.Add(item.SellOrders[0].RemainQuantity)
				deals = append(deals, Deal{item.SellOrders[0].OrderId, "SELL", item.SellOrders[0].RemainQuantity})
				item.SellOrders = item.SellOrders[1:]
			} else {
				deals = append(deals, Deal{item.SellOrders[0].OrderId, "SELL", sellQuantity.Sub(filledSellQuantity)})
				item.SellOrders[0].RemainQuantity = item.SellOrders[0].RemainQuantity.Sub(sellQuantity.Sub(filledSellQuantity))
				break
			}
		}
	}
	return deals
}

func PeriodicAuctionMatch(book *DepthBook, pricePrecision int64, refPrice sdk.Dec) *MatchResult {
	if book == nil {
		return &MatchResult{sdk.ZeroDec(), sdk.ZeroDec(), []Deal{}}
	}
	bestPrice, maxExecution := periodicAuctionMatchPrice(book, pricePrecision, refPrice)
	if maxExecution.IsZero() {
		return &MatchResult{sdk.ZeroDec(), sdk.ZeroDec(), []Deal{}}
	}
	buyAmount := sdk.ZeroDec()
	sellAmount := sdk.ZeroDec()
	deals := []Deal{}

	index := 0
	bestPriceIndex := len(*book)
	for index < len(*book) {
		if (*book)[index].Price.GT(bestPrice) { // item.Price > bestPrice
			buyAmount = buyAmount.Add((*book)[index].BuyQuantity)
			deals = append(deals, fillDepthBookItem(&(*book)[index], (*book)[index].BuyQuantity, sdk.ZeroDec())...)
			if (*book)[index].SellQuantity.IsZero() { // this item has no buy quantity or sell quantity anymore
				*book = append((*book)[:index], (*book)[index+1:]...)
			} else {
				index++
			}
		} else if (*book)[index].Price.LT(bestPrice) { // item.Price < bestPrice
			sellAmount = sellAmount.Add((*book)[index].SellQuantity)
			deals = append(deals, fillDepthBookItem(&(*book)[index], sdk.ZeroDec(), (*book)[index].SellQuantity)...)
			if (*book)[index].BuyQuantity.IsZero() { // this item has no buy quantity or sell quantity anymore
				*book = append((*book)[:index], (*book)[index+1:]...)
			} else {
				index++
			}
		} else { // item.Price == bestPrice
			bestPriceIndex = index
			index++
		}
	}
	if bestPriceIndex < len(*book) {
		deals = append(deals, fillDepthBookItem(&(*book)[bestPriceIndex], maxExecution.Sub(buyAmount), maxExecution.Sub(sellAmount))...)
	}
	return &MatchResult{
		bestPrice,
		maxExecution,
		deals,
	}
}
