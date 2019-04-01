package order

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgNewOrder:
			return handleMsgNewOrder(ctx, keeper, msg)
		case MsgCancelOrder:
			return handleMsgCancelOrder(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized order Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgNewOrder
func handleMsgNewOrder(ctx sdk.Context, keeper Keeper, msg MsgNewOrder) sdk.Result {
	price, _ := sdk.NewDecFromStr(msg.Price)
	quantity, _ := sdk.NewDecFromStr(msg.Quantity)
	order := &Order{
		Sender:         msg.Sender,
		Product:        msg.Product,
		Side:           msg.Side,
		Price:          price,
		Quantity:       quantity,
		Status:         OrderStatusOpen,
		RemainQuantity: quantity,
	}
	// Trying to lock coins
	needLockCoins := order.NeedLockCoins()
	err := keeper.LockCoins(ctx, msg.Sender, needLockCoins)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInsufficientCoins,
		}
	}
	blockHeight := ctx.BlockHeight()
	orderNum := keeper.GetBlockOrderNum(ctx, blockHeight)
	order.OrderId = FormatOrderId(blockHeight, orderNum)
	keeper.SetBlockOrderNum(ctx, blockHeight, orderNum+1)
	keeper.SetOrder(ctx, order.OrderId, order)

	return sdk.Result{
		Tags: sdk.NewTags(TagKeyOrderId, order.OrderId),
	}
}

func handleMsgCancelOrder(ctx sdk.Context, keeper Keeper, msg MsgCancelOrder) sdk.Result {
	order := keeper.GetOrder(ctx, msg.OrderId)
	if order.Sender.String() != msg.Sender.String() {
		return sdk.Result{
			Code: sdk.CodeUnauthorized,
			Log:  fmt.Sprintf("not the owner of order(%v)", msg.OrderId),
		}
	}
	if order.Status != OrderStatusOpen {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("cannot cancel order with status(%v)", order.Status),
		}
	}
	order.Cancel()
	// unlock coins in this order & charge fee
	needUnlockCoins := order.NeedUnlockCoins()
	fee, inOrder := GetOrderCancelFee(order, ctx, keeper)
	keeper.AddCollectedFees(ctx, fee)
	if inOrder {
		needUnlockCoins = needUnlockCoins.Sub(fee)
		keeper.BurnLockedCoins(ctx, order.Sender, fee)
	} else {
		keeper.SubtractCoins(ctx, order.Sender, fee)
	}
	keeper.UnlockCoins(ctx, msg.Sender, needUnlockCoins)
	keeper.SetOrder(ctx, order.OrderId, order)

	// update depth book
	depthBook := keeper.GetDepthBook(ctx, order.Product)
	depthBook.RemoveOrder(order)
	keeper.SetDepthBook(ctx, order.Product, depthBook)

	return sdk.Result{
		Tags: sdk.NewTags(TagKeyOrderId, order.OrderId),
	}
}
