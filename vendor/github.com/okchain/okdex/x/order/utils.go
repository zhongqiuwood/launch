package order

// TODO: move to common

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ConvertDecCoinsToCoins(decCoins sdk.DecCoins) sdk.Coins {
	coins := sdk.Coins{}
	for _, decCoin := range decCoins {
		coin := sdk.NewCoin(decCoin.Denom, sdk.NewIntFromBigInt(decCoin.Amount.Int))
		coins = append(coins, coin)
	}
	coins = coins.Sort()
	return coins
}

func ConvertCoinsToDecCoins(coins sdk.Coins) sdk.DecCoins {
	decCoins := sdk.DecCoins{}
	for _, coin := range coins {
		decCoin := sdk.NewDecCoinFromDec(coin.Denom, sdk.NewDecFromIntWithPrec(coin.Amount, sdk.Precision))
		decCoins = append(decCoins, decCoin)
	}
	decCoins = decCoins.Sort()
	return decCoins
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}
