package order

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// TODO: get params from param keeper
const (
	FeeCancel         = "0.01"
	FeeCancelNative   = "0.002"
	FeeExpire         = "0.01"
	FeeExpireNative   = "0.002"
	FeeRateDeal       = "0.001"
	FeeRateDealNative = "0.0004"
)

// Cancel/Expire fees are charged according to the priority below:
// 1. native okb locked in order
// 2. native okb in account balance
// 3. other token locked in order
func GetOrderCancelFee(order *Order, ctx sdk.Context, keeper Keeper) (fee sdk.DecCoins, inOrder bool) {
	// 1. native okb locked in order
	symbols := strings.Split(order.Product, "_")
	feeNativeAmt := sdk.MustNewDecFromStr(FeeCancelNative)
	// Make sure feeNativeAmt <= remain value in order
	feeNativeAmt = sdk.MinDec(feeNativeAmt, order.Price.Mul(order.RemainQuantity))
	feeNative := sdk.DecCoins{sdk.NewDecCoinFromDec("okb", feeNativeAmt)}
	if (symbols[1] == "okb" && order.Side == "BUY") || (symbols[0] == "okb" && order.Side == "SELL") {
		return feeNative, true
	}
	// 2. native okb in account balance
	if keeper.HasCoins(ctx, order.Sender, feeNative) {
		return feeNative, false
	}
	// 3. other token locked in order
	symbol := symbols[0]
	if order.Side == "BUY" {
		symbol = symbols[1]
	}
	product := fmt.Sprintf("%s_okb", symbol)
	lastPrice := keeper.GetLastPrice(ctx, product)
	feeOtherAmt := sdk.MustNewDecFromStr(FeeCancel).Quo(lastPrice)
	// Make sure feeOtherAmt <= remain quantity in order
	if order.Side == "SELL" {
		feeOtherAmt = sdk.MinDec(feeOtherAmt, order.RemainQuantity)
	} else {
		feeOtherAmt = sdk.MinDec(feeOtherAmt, order.RemainQuantity.Mul(order.Price))
	}
	feeOther := sdk.DecCoins{sdk.NewDecCoinFromDec(symbol, feeOtherAmt)}
	return feeOther, true
}

func GetOrderExpireFee(order *Order, ctx sdk.Context, keeper Keeper) (fee sdk.DecCoins, inOrder bool) {
	// 1. native okb locked in order
	symbols := strings.Split(order.Product, "_")
	feeNativeAmt := sdk.MustNewDecFromStr(FeeExpireNative)
	// Make sure feeNativeAmt <= remain value of order
	feeNativeAmt = sdk.MinDec(feeNativeAmt, order.Price.Mul(order.RemainQuantity))
	feeNative := sdk.DecCoins{sdk.NewDecCoinFromDec("okb", feeNativeAmt)}
	if (symbols[1] == "okb" && order.Side == "BUY") || (symbols[0] == "okb" && order.Side == "SELL") {
		return feeNative, true
	}
	// 2. native okb in account balance
	if keeper.HasCoins(ctx, order.Sender, feeNative) {
		return feeNative, false
	}
	// 3. other token locked in order
	symbol := symbols[0]
	if order.Side == "BUY" {
		symbol = symbols[1]
	}
	product := fmt.Sprintf("%s_okb", symbol)
	lastPrice := keeper.GetLastPrice(ctx, product)
	feeOther := sdk.DecCoins{sdk.NewDecCoinFromDec(symbol, sdk.MustNewDecFromStr(FeeExpire).Quo(lastPrice))}
	return feeOther, true
}

// Deal fees are charged according to the priority below:
// 1. native okb in account balance
// 2. native okb received in order
// 3. other token received in order
func GetDealFee(deal *Deal, order *Order, ctx sdk.Context, keeper Keeper) (fee sdk.DecCoins, inOrder bool) {
	// 1. native okb in account balance
	symbols := strings.Split(order.Product, "_")
	okbAmt := deal.Quantity.Mul(keeper.GetLastPrice(ctx, order.Product))
	if symbols[1] != "okb" {
		productQuote := fmt.Sprintf("%s_okb", symbols[1])
		okbAmt = okbAmt.Mul(keeper.GetLastPrice(ctx, productQuote))
	}
	feeAmt := okbAmt.Mul(sdk.MustNewDecFromStr(FeeRateDealNative))
	feeNative := sdk.DecCoins{sdk.NewDecCoinFromDec("okb", feeAmt)}
	if keeper.HasCoins(ctx, order.Sender, feeNative) {
		return feeNative, false
	}

	// 2. native okb received in order
	if symbols[1] == "okb" && order.Side == "SELL" {
		return feeNative, true
	}

	// 3. other token received in order
	symbol := symbols[0]
	quantity := deal.Quantity
	if order.Side == "SELL" {
		symbol = symbols[1]
		quantity = deal.Quantity.Mul(keeper.GetLastPrice(ctx, order.Product))
	}
	feeOther := sdk.DecCoins{sdk.NewDecCoinFromDec(symbol, quantity.Mul(sdk.MustNewDecFromStr(FeeRateDeal)))}
	return feeOther, true
}
