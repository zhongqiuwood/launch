package common

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	EPOCH_BLOCKS = int64(3 * 24 * 60 * 60) //3days, 1s per block
)

//judge if reach one epoch end, if not, return the remaining interval
func IsEpochEnd(ctx sdk.Context) (remainder int64, ok bool) {
	if ctx.BlockHeight()%GetEpochInterval() == 0 {
		return 0, true
	}
	return GetEpochInterval() - ctx.BlockHeight()%GetEpochInterval(), false
}

func GetEpochInterval() int64 {
	return EPOCH_BLOCKS
}
