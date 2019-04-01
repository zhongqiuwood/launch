package backend

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Called every block, check expired orders
func EndBlocker(ctx sdk.Context, keeper Keeper, orderKeeper OrderKeeper) sdk.Tags {
	resTags := storeTradeAndKLine(ctx, keeper, orderKeeper)
	return resTags
}

func storeTradeAndKLine(ctx sdk.Context, keeper Keeper, orderKeeper OrderKeeper) sdk.Tags {
	resTags := sdk.NewTags()
	blockHeight := ctx.BlockHeight()
	timestamp := ctx.BlockHeader().Time.Unix()
	matchResultMap := orderKeeper.GetMatchResultMap(ctx, ctx.BlockHeight())
	for product, matchResult := range matchResultMap {
		match := &Match{
			BlockHeight: blockHeight,
			Product:     product,
			Price:       matchResult.Price.String(),
			Quantity:    matchResult.Quantity.String(),
			Timestamp:   timestamp,
		}
		keeper.StoreMatch(ctx, match)
		for _, record := range matchResult.Deals {
			order := orderKeeper.GetOrder(ctx, record.OrderId)
			trade := &Trade{
				BlockHeight: blockHeight,
				OrderId:     record.OrderId,
				Sender:      order.Sender.String(),
				Product:     product,
				Price:       matchResult.Price.String(),
				Quantity:    record.Quantity.String(),
				Timestamp:   timestamp,
			}
			keeper.StoreTrade(ctx, trade)
		}
		// TODO: compute kline and store
	}
	return resTags
}
