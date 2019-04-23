package order

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func handleMsgNewOrderV2(ctx sdk.Context, keeper Keeper, msg MsgNewOrder) sdk.Result {
	return handleMsgNewOrder(ctx, keeper, msg)
}

func handleMsgCancelOrderV2(ctx sdk.Context, keeper Keeper, msg MsgCancelOrder) sdk.Result {
	return handleMsgCancelOrder(ctx, keeper, msg)
}
