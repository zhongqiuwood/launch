package backend

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Called every block, check expired orders
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	if keeper.logger == nil {
		keeper.logger = ctx.Logger().With("/x/backend")
	}

	if keeper.maintainConf.EnableBackend {
		event := &EndBlockEvent{ctx, ctx.BlockHeight(), ctx.BlockHeader().Time.Unix()}
		keeper.endBlockerChan <- event
	}
}
