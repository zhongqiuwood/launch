package backend

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Called every block, check expired orders
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	if keeper.maintainConf.EnableBackend {
		event := &EndBlockEvent{ctx, ctx.BlockHeight(), ctx.BlockHeader().Time.Unix()}
		keeper.ProductEndBlockerEvent(event)
	}
}
