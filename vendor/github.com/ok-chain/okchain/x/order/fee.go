package order

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ok-chain/okchain/x/common"
	"strings"
)

type GetFeeKeeper interface {
	HasCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins) bool
	GetLastPrice(ctx sdk.Context, product string) sdk.Dec
}

// Currently, placing order does not need any fee, so we only support charging okb if necessary
func GetOrderNewFee(order *Order, ctx sdk.Context, keeper GetFeeKeeper, feeParams Params) (sdk.DecCoins, bool) {
	return sdk.DecCoins{sdk.NewDecCoinFromDec(common.ChainAsset, feeParams.NewOrder)}, false
}

// Cancel/Expire fees are charged according to the priority below:
// 1. native okb in account balance
// 2. other token locked in order
func GetOrderCancelFee(order *Order, ctx sdk.Context, keeper GetFeeKeeper, feeParams Params) sdk.DecCoins {
	// 1. native okb in account balance
	symbols := strings.Split(order.Product, "_")
	feeNativeAmt := feeParams.CancelNative
	feeNative := sdk.DecCoins{sdk.NewDecCoinFromDec(common.ChainAsset, feeNativeAmt)}
	if keeper.HasCoins(ctx, order.Sender, feeNative) {
		return feeNative
	}

	// 2. other token locked in order
	symbol := symbols[0]
	if order.Side == "BUY" {
		symbol = symbols[1]
	}
	product := fmt.Sprintf("%s_%s", symbol, common.ChainAsset)
	lastPrice := keeper.GetLastPrice(ctx, product)
	feeOtherAmt := feeParams.Cancel.Quo(lastPrice)
	// Make sure feeOtherAmt <= remain quantity in order
	if order.Side == "SELL" {
		feeOtherAmt = sdk.MinDec(feeOtherAmt, order.RemainQuantity)
	} else {
		feeOtherAmt = sdk.MinDec(feeOtherAmt, order.RemainQuantity.Mul(order.Price))
	}
	feeOther := sdk.DecCoins{sdk.NewDecCoinFromDec(symbol, feeOtherAmt)}
	return feeOther
}

func GetOrderExpireFee(order *Order, ctx sdk.Context, keeper GetFeeKeeper, feeParams Params) sdk.DecCoins {
	// 1. native okb in account balance
	symbols := strings.Split(order.Product, "_")
	feeNativeAmt := feeParams.ExpireNative
	feeNative := sdk.DecCoins{sdk.NewDecCoinFromDec(common.ChainAsset, feeNativeAmt)}
	if keeper.HasCoins(ctx, order.Sender, feeNative) {
		return feeNative
	}

	// 2. other token locked in order
	symbol := symbols[0]
	if order.Side == "BUY" {
		symbol = symbols[1]
	}
	product := fmt.Sprintf("%s_%s", symbol, common.ChainAsset)
	lastPrice := keeper.GetLastPrice(ctx, product)
	feeOtherAmt := feeParams.Expire.Quo(lastPrice)
	// Make sure feeOtherAmt <= remain quantity in order
	if order.Side == "SELL" {
		feeOtherAmt = sdk.MinDec(feeOtherAmt, order.RemainQuantity)
	} else {
		feeOtherAmt = sdk.MinDec(feeOtherAmt, order.RemainQuantity.Mul(order.Price))
	}
	feeOther := sdk.DecCoins{sdk.NewDecCoinFromDec(symbol, feeOtherAmt)}
	return feeOther
}

// Deal fees are charged according to the priority below:
// 1. native okb in account balance
// 2. other token received in order
func GetDealFee(order *Order, fillAmt sdk.Dec, ctx sdk.Context, keeper GetFeeKeeper, feeParams Params) sdk.DecCoins {
	// 1. native okb in account balance
	symbols := strings.Split(order.Product, "_")
	okbAmt := fillAmt.Mul(keeper.GetLastPrice(ctx, order.Product))
	if symbols[1] != common.ChainAsset {
		productQuote := fmt.Sprintf("%s_%s", symbols[1], common.ChainAsset)
		okbAmt = okbAmt.Mul(keeper.GetLastPrice(ctx, productQuote))
	}
	feeAmt := okbAmt.Mul(feeParams.TradeFeeRateNative)
	feeNative := sdk.DecCoins{sdk.NewDecCoinFromDec(common.ChainAsset, feeAmt)}
	if keeper.HasCoins(ctx, order.Sender, feeNative) {
		return feeNative
	}

	// 2. other token received in order
	symbol := symbols[0]
	quantity := fillAmt
	if order.Side == "SELL" {
		symbol = symbols[1]
		quantity = fillAmt.Mul(keeper.GetLastPrice(ctx, order.Product))
	}
	feeOther := sdk.DecCoins{sdk.NewDecCoinFromDec(symbol, quantity.Mul(feeParams.TradeFeeRate))}
	return feeOther
}
