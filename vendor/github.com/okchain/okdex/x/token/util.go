package token

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ToUnit(amount int64) sdk.Int {
	return sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(amount), new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))
}

func FloatToCoins(fee string) sdk.Coins {
	feeDec := sdk.MustNewDecFromStr(fee)
	coin := sdk.NewCoin("okb", sdk.NewIntFromBigInt(feeDec.Int))
	var coins sdk.Coins
	coins = append(coins, coin)
	return coins
}
