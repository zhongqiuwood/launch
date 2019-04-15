package order

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
)

// Round a decimal with precision, perform bankers rounding (gaussian rounding)
func RoundDecimal(dec sdk.Dec, precision int64) sdk.Dec {
	precisionMul := sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(precision), nil))
	return sdk.NewDecFromInt(dec.MulInt(precisionMul).RoundInt()).QuoInt(precisionMul)
}

// Calculate periodic auction match price, return the best price and execution amount
// The best price is found according following rules:
// rule0: No match, bestPrice = 0, maxExecution=0
// rule1: Maximum execution volume. If there are more than one price with the same max execution, following rule2
// rule2: Minimum imbalance. We should select the price with minimum absolute value of imbalance. If more
//        than one price satisfy rule2, following rule3
// rule3: Market Pressure. There are 3 cases:
// rule3a: All imbalances are positive. It indicates buy side pressure. Set reference price as last execute price
//         plus a upper limit percentage(e.g. 5%). Then choose the price which is closest to reference price.
// rule3b: All imbalances are negative. It indicates sell side pressure. Set reference price as last execute price
//         minus a lower limit percentage(e.g. 5%). Then choose the price which is closest to reference price.
// rule3a: Otherwise, it indicates no one side pressure. Set reference price as last execute price.
//         Then choose the price which is closest to reference price.
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

	// See rule 0: no match
	if maxExecution.IsZero() {
		return refPrice, maxExecution
	}

	var indexesRule1 []int  // price indexes satisfied rule1
	for i := 0; i < bookLength; i++ {
		if execution[i].Equal(maxExecution) {
			indexesRule1 = append(indexesRule1, i)
		}
	}
	indexLen1 := len(indexesRule1)
	// See rule1: Maximum execution
	if indexLen1 == 1 {
		bestPrice = (*book)[indexesRule1[0]].Price
		return
	}

	// See rule2: Minimum imbalance
	imbalance := make([]sdk.Dec, bookLength)
	for i := 0; i < bookLength; i++ {
		imbalance[i] = buyAmountSum[i].Sub(sellAmountSum[i])
	}
	minAbsImbalance := imbalance[indexesRule1[0]].Abs()
	for i := 1; i < indexLen1; i++ {
		minAbsImbalance = sdk.MinDec(minAbsImbalance, imbalance[indexesRule1[i]].Abs())
	}
	var indexesRule2 []int  // price indexes satisfied rule1&rule2
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

	// See rule3: Market Pressure
	if imbalance[indexesRule2[0]].GT(sdk.ZeroDec()) { // rule3a: all imbalances are positive, buy side pressure
		newRefPrice := refPrice.Mul(sdk.MustNewDecFromStr("1.05"))
		newRefPrice = RoundDecimal(newRefPrice, pricePrecision)
		if (*book)[indexesRule2[0]].Price.LT(newRefPrice) { // all price < newRefPrice, choose max price
			bestPrice = (*book)[indexesRule2[0]].Price
		} else if (*book)[indexesRule2[indexLen2-1]].Price.GT(newRefPrice) { // all price > newRefPrice, choose min price
			bestPrice = (*book)[indexesRule2[indexLen2-1]].Price
		} else {  // otherwise, choose newRefPrice
			bestPrice = newRefPrice
		}
	} else if imbalance[indexesRule2[indexLen2-1]].LT(sdk.ZeroDec()) { // rule3b: all imbalances are negative, sell side pressure
		newRefPrice := refPrice.Mul(sdk.MustNewDecFromStr("0.95"))
		newRefPrice = RoundDecimal(newRefPrice, pricePrecision)
		if (*book)[indexesRule2[0]].Price.LT(newRefPrice) { // all price < newRefPrice, choose max price
			bestPrice = (*book)[indexesRule2[0]].Price
		} else if (*book)[indexesRule2[indexLen2-1]].Price.GT(newRefPrice) { // all price > newRefPrice, choose min price
			bestPrice = (*book)[indexesRule2[indexLen2-1]].Price
		} else {  // otherwise, choose newRefPrice
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
