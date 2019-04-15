package order

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type GetFeeKeeper interface {
	HasCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins) bool
	GetLastPrice(ctx sdk.Context, product string) sdk.Dec
}

// Currently, placing order does not need any fee, so we only support charging okb if necessary
func GetOrderNewFee(order *Order, ctx sdk.Context, keeper GetFeeKeeper, feeParams Params) sdk.DecCoins {
	if feeParams.NewOrder.IsZero() {
		return sdk.DecCoins{}
	}
	return sdk.DecCoins{sdk.NewDecCoinFromDec("okb", feeParams.NewOrder)}
}

// Cancel/Expire fees are charged according to the priority below:
// 1. native okb locked in order
// 2. native okb in account balance
// 3. other token locked in order
func GetOrderCancelFee(order *Order, ctx sdk.Context, keeper GetFeeKeeper, feeParams Params) (fee sdk.DecCoins, inOrder bool) {
	// 1. native okb locked in order
	symbols := strings.Split(order.Product, "_")
	feeNativeAmt := feeParams.CancelNative
	// Make sure feeNativeAmt <= remain value in order
	feeNativeAmt = sdk.MinDec(feeNativeAmt, order.Price.Mul(order.RemainQuantity))
	feeNative := sdk.DecCoins{sdk.NewDecCoinFromDec("okb", feeNativeAmt)}
	if symbols[1] == "okb" && order.Side == "BUY" {
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
	feeOtherAmt := feeParams.Cancel.Quo(lastPrice)
	// Make sure feeOtherAmt <= remain quantity in order
	if order.Side == "SELL" {
		feeOtherAmt = sdk.MinDec(feeOtherAmt, order.RemainQuantity)
	} else {
		feeOtherAmt = sdk.MinDec(feeOtherAmt, order.RemainQuantity.Mul(order.Price))
	}
	feeOther := sdk.DecCoins{sdk.NewDecCoinFromDec(symbol, feeOtherAmt)}
	return feeOther, true
}

func GetOrderExpireFee(order *Order, ctx sdk.Context, keeper GetFeeKeeper, feeParams Params) (fee sdk.DecCoins, inOrder bool) {
	// 1. native okb locked in order
	symbols := strings.Split(order.Product, "_")
	feeNativeAmt := feeParams.ExpireNative
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
	feeOtherAmt := feeParams.Expire.Quo(lastPrice)
	// Make sure feeOtherAmt <= remain quantity in order
	if order.Side == "SELL" {
		feeOtherAmt = sdk.MinDec(feeOtherAmt, order.RemainQuantity)
	} else {
		feeOtherAmt = sdk.MinDec(feeOtherAmt, order.RemainQuantity.Mul(order.Price))
	}
	feeOther := sdk.DecCoins{sdk.NewDecCoinFromDec(symbol, feeOtherAmt)}
	return feeOther, true
}

// Deal fees are charged according to the priority below:
// 1. native okb in account balance
// 2. native okb received in order
// 3. other token received in order
func GetDealFee(order *Order, fillAmt sdk.Dec, ctx sdk.Context, keeper GetFeeKeeper, feeParams Params) (fee sdk.DecCoins, inOrder bool) {
	// 1. native okb in account balance
	symbols := strings.Split(order.Product, "_")
	okbAmt := fillAmt.Mul(keeper.GetLastPrice(ctx, order.Product))
	if symbols[1] != "okb" {
		productQuote := fmt.Sprintf("%s_okb", symbols[1])
		okbAmt = okbAmt.Mul(keeper.GetLastPrice(ctx, productQuote))
	}
	feeAmt := okbAmt.Mul(feeParams.TradeFeeRateNative)
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
	quantity := fillAmt
	if order.Side == "SELL" {
		symbol = symbols[1]
		quantity = fillAmt.Mul(keeper.GetLastPrice(ctx, order.Product))
	}
	feeOther := sdk.DecCoins{sdk.NewDecCoinFromDec(symbol, quantity.Mul(feeParams.TradeFeeRate))}
	return feeOther, true
}
