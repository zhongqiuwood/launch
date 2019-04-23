package order

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func endBlockerV2(ctx sdk.Context, keeper Keeper) sdk.Tags {
	return endBlockerV1(ctx, keeper)
}

