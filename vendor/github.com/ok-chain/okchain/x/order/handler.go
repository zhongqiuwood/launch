package order

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ok-chain/okchain/x/common"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"strconv"
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

var mockBlockHeight int64 = -1

func blockHeight(ctx sdk.Context) int64 {

	if mockBlockHeight >= 0 {
		return mockBlockHeight
	}

	return ctx.BlockHeight()
}

// checkOrder: check order product, price & quantity fields
func checkOrder(ctx sdk.Context, keeper Keeper, order *Order) error {
	tokenPair := keeper.tokenKeeper.GetTokenPair(ctx, order.Product)
	if tokenPair == nil {
		return errors.New(fmt.Sprintf("Invalid product: %s", order.Product))
	}

	priceDigit := int64(tokenPair.MaxPriceDigit)
	quantityDigit := int64(tokenPair.MaxQuantityDigit)
	roundedPrice := RoundDecimal(order.Price, priceDigit)
	roundedQuantity := RoundDecimal(order.Quantity, quantityDigit)
	if !roundedPrice.Equal(order.Price) {
		return errors.New(fmt.Sprintf("Price(%v) over accuracy(%d)", order.Price, priceDigit))
	}
	if !roundedQuantity.Equal(order.Quantity) {
		return errors.New(fmt.Sprintf("Quantity(%v) over accuracy(%d)", order.Quantity, quantityDigit))
	}

	if order.Quantity.LT(tokenPair.MinQuantity) {
		return errors.New(fmt.Sprintf("Quantity should be greater than %s", tokenPair.MinQuantity))
	}
	return nil
}

func placeOrder(ctx sdk.Context, keeper Keeper, order *Order, feeParams Params) error {
	if !feeParams.NewOrder.IsZero() {
		fee := GetOrderNewFee(order, ctx, keeper, feeParams)
		err := keeper.SubtractCoins(ctx, order.Sender, fee)
		if err != nil {
			keeper.AddCollectedFees(ctx, fee, order.Sender, common.FeeTypeOrderNew)
			return err
		}
	}

	// Trying to lock coins
	needLockCoins := order.NeedLockCoins()
	err := keeper.LockCoins(ctx, order.Sender, needLockCoins)
	if err != nil {
		dump(ctx, "%s! order[%v]", err, order)
		return err
	}
	blockHeight := blockHeight(ctx)
	orderNum := keeper.GetBlockOrderNum(ctx, blockHeight)
	order.OrderId = common.FormatOrderId(blockHeight, orderNum)

	dump(ctx, "new order[%v]", order)
	keeper.SetBlockOrderNum(ctx, blockHeight, orderNum+1)
	keeper.SetOrder(ctx, order.OrderId, order)
	return nil
}

// Handle MsgNewOrder
func handleMsgNewOrder(ctx sdk.Context, keeper Keeper, msg MsgNewOrder) sdk.Result {
	price, _ := sdk.NewDecFromStr(msg.Price)
	quantity, _ := sdk.NewDecFromStr(msg.Quantity)
	order := &Order{
		TxHash:         fmt.Sprintf("%X", tmhash.Sum(ctx.TxBytes())),
		Sender:         msg.Sender,
		Product:        msg.Product,
		Side:           msg.Side,
		Price:          price,
		Quantity:       quantity,
		Status:         OrderStatusOpen,
		RemainQuantity: quantity,
		Timestamp:      ctx.BlockHeader().Time.Unix(),
	}

	err := checkOrder(ctx, keeper, order)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeUnknownRequest,
			Log:  err.Error(),
		}
	}

	//just for test
	batch_number := 1
	num_t, err := strconv.ParseInt(msg.BatchNumber, 10, 64)
	if err == nil && num_t > 0 {
		batch_number = int(num_t)
	}
	const max_batch_number = 10000000
	if batch_number < 1 {
		batch_number = 1
	}
	if batch_number > max_batch_number {
		batch_number = max_batch_number
	}
	var tags sdk.Tags
	for i := 0; i < batch_number; i++ {
		feeParams := keeper.GetParams(ctx)
		if err := placeOrder(ctx, keeper, order, feeParams); err != nil {
			return sdk.Result{
				Code: sdk.CodeInsufficientCoins,
				Log: err.Error(),
			}
		}
		if i == 0 {
			tags = append(tags, sdk.NewTags(TagKeyOrderId, order.OrderId)...)
		}

	}
	tags = append(tags, sdk.NewTags("batch_number", strconv.Itoa(batch_number))...)
	return sdk.Result{
		Tags: tags,
	}
	//return sdk.Result{
	//	Tags: sdk.NewTags(TagKeyOrderId, order.OrderId),
	//}
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
	if order.Status == OrderStatusCancelled { // charge fees only if fully cancelled
		feeParams := keeper.GetParams(ctx)
		fee, inOrder := GetOrderCancelFee(order, ctx, keeper, feeParams)
		keeper.AddCollectedFees(ctx, fee, order.Sender, common.FeeTypeOrderCancel)
		if inOrder {
			needUnlockCoins = needUnlockCoins.Sub(fee)
			keeper.BurnLockedCoins(ctx, order.Sender, fee)
		} else {
			keeper.SubtractCoins(ctx, order.Sender, fee)
		}
	}
	keeper.UnlockCoins(ctx, msg.Sender, needUnlockCoins)
	keeper.SetOrder(ctx, order.OrderId, order)

	// update depth book and orderIdsMap
	orderIdsMap := make(OrderIdsMap)
	depthBook := keeper.GetDepthBook(ctx, order.Product)
	depthBook.RemoveOrder(order)
	orderIdsMap.RemoveOrder(order, ctx, keeper)
	keeper.SetDepthBook(ctx, order.Product, depthBook)
	keeper.AddUpdatedOrderId(ctx, ctx.BlockHeight(), order.OrderId)
	orderIdsMap.SaveToKeeper(ctx, keeper)

	return sdk.Result{
		Tags: sdk.NewTags(TagKeyOrderId, order.OrderId),
	}
}
