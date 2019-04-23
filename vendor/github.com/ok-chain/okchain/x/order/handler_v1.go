package order

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ok-chain/okchain/x/common"
	"github.com/ok-chain/okchain/x/perf"
	"github.com/ok-chain/okchain/x/version"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"strconv"
)

func NewHandler(keeper Keeper) sdk.Handler {
	handler := func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		return version.GetVersionController().NewHandler(ctx, msg, keeper, getModule())
	}
	return handler
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
	// charge fee for placing a new order
	if !feeParams.NewOrder.IsZero() {
		fee, _ := GetOrderNewFee(order, ctx, keeper, feeParams)
		if err := keeper.SubtractCoins(ctx, order.Sender, fee); err == nil {
			keeper.AddCollectedFees(ctx, fee, order.Sender, common.FeeTypeOrderNew)
			order.RecordOrderNewFee(fee)
		} else {
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

func PlaceOrder(ctx sdk.Context, keeper Keeper, order *Order, feeParams Params) error {
	return placeOrder(ctx, keeper, order, feeParams)
}

// Handle MsgNewOrder
func handleMsgNewOrder(ctx sdk.Context, keeper Keeper, msg MsgNewOrder) sdk.Result {

	seq := perf.GetPerf().OnDeliverTxEnter(ctx, ModuleName, "handleMsgNewOrder")
	defer perf.GetPerf().OnDeliverTxExit(ctx, ModuleName, "handleMsgNewOrder", seq)

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
				Log:  err.Error(),
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

	seq := perf.GetPerf().OnDeliverTxEnter(ctx, ModuleName, "handleMsgCancelOrder")
	defer perf.GetPerf().OnDeliverTxExit(ctx, ModuleName, "handleMsgCancelOrder", seq)

	order := keeper.GetOrder(ctx, msg.OrderId)
	feeParams := keeper.GetParams(ctx)

	// If order is invalid, try to charge cancel fee as penalty
	if order == nil || order.Sender.String() != msg.Sender.String() || order.Status != OrderStatusOpen {
		feePenalty := sdk.DecCoins{sdk.NewDecCoinFromDec(common.ChainAsset, feeParams.CancelNative)}
		if err := keeper.SubtractCoins(ctx, msg.Sender, feePenalty); err == nil {
			keeper.AddCollectedFees(ctx, feePenalty, msg.Sender, common.FeeTypeOrderCancel)
		}

		if order == nil || order.Status != OrderStatusOpen {
			return sdk.Result{
				Code: sdk.CodeInternal,
				Log:  fmt.Sprintf("cannot cancel order(%v)", order),
			}
		} else {
			return sdk.Result{
				Code: sdk.CodeUnauthorized,
				Log:  fmt.Sprintf("not the owner of order(%v)", msg.OrderId),
			}
		}
	}

	order.Cancel()
	// unlock coins in this order & charge fee
	needUnlockCoins := order.NeedUnlockCoins()
	keeper.UnlockCoins(ctx, msg.Sender, needUnlockCoins)
	if order.Status == OrderStatusCancelled { // charge fees only if fully cancelled
		fee := GetOrderCancelFee(order, ctx, keeper, feeParams)
		keeper.SubtractCoins(ctx, order.Sender, fee)
		keeper.AddCollectedFees(ctx, fee, order.Sender, common.FeeTypeOrderCancel)
		order.RecordOrderCancelFee(fee)
	}
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
