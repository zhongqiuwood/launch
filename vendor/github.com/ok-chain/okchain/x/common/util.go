package common

import (
	"encoding/binary"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"runtime/debug"
)

// ConvertDecCoinsToCoins return coins by multiplying decCoins times 1e18, eg. 0.000000001okb -> 1000000000okb
func ConvertDecCoinsToCoins(decCoins sdk.DecCoins) sdk.Coins {
	coins := sdk.Coins{}
	for _, decCoin := range decCoins {
		coin := sdk.NewCoin(decCoin.Denom, sdk.NewIntFromBigInt(decCoin.Amount.Int))
		coins = append(coins, coin)
	}
	coins = coins.Sort()
	return coins
}

// ConvertCoinsToDecCoins return decCoins by dividing coins by 1e18, eg. 1000000000okb -> 0.000000001okb
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

func GetPage(page, perPage int) (offset, limit int) {
	offset = (page - 1) * perPage
	limit = perPage
	return
}

func FormatOrderId(blockHeight, orderNum int64) string {
	format := "ID%010d-%d"
	if blockHeight > 9999999999 {
		format = "ID%d-%d"
	}
	return fmt.Sprintf(format, blockHeight, orderNum)
}

func PrintStackIfPainic()  {
	r := recover()
	if r != nil {
		debug.PrintStack()
	}
}