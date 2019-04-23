package order

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ok-chain/okchain/x/version"
	"sync"
)

type orderModule struct {}
var module *orderModule
var once sync.Once

func getModule() version.AppModule {
	once.Do(func() {
		module = &orderModule{}
	})
	return module
}


func (m *orderModule) NewHandlerVersion1(ctx sdk.Context, msg sdk.Msg, k interface{}) sdk.Result {

	keeper, ok := k.(Keeper)
	if !ok {
		return m.invalidKeeperResult()
	}

	switch msg := msg.(type) {
	case MsgNewOrder:
		return handleMsgNewOrder(ctx, keeper, msg)
	case MsgCancelOrder:
		return handleMsgCancelOrder(ctx, keeper, msg)
	default:
		return m.invalidMessageResult(msg)
	}
}

func (m *orderModule) NewHandlerVersion2(ctx sdk.Context, msg sdk.Msg, k interface{}) sdk.Result {

	keeper, ok := k.(Keeper)
	if !ok {
		return m.invalidKeeperResult()
	}

	switch msg := msg.(type) {
	case MsgNewOrder:
		return handleMsgNewOrderV2(ctx, keeper, msg)
	case MsgCancelOrder:
		return handleMsgCancelOrderV2(ctx, keeper, msg)
	default:
		return m.invalidMessageResult(msg)
	}
}


func (m *orderModule) EndBlockerVersion1(ctx sdk.Context, k interface{}) sdk.Tags {

	keeper, ok := k.(Keeper)

	if !ok {
		return m.invalidKeeperTags()
	}

	return endBlockerV1(ctx, keeper)
}

func (m *orderModule) EndBlockerVersion2(ctx sdk.Context, k interface{}) sdk.Tags {

	keeper, ok := k.(Keeper)

	if !ok {
		return m.invalidKeeperTags()
	}

	return endBlockerV2(ctx, keeper)
}

func (m *orderModule) invalidMessageResult(msg sdk.Msg) sdk.Result {
	errMsg := fmt.Sprintf("Unrecognized order Msg type: %v", msg.Type())
	return sdk.ErrUnknownRequest(errMsg).Result()
}

func (m *orderModule) invalidKeeperResult() sdk.Result {
	errMsg := fmt.Sprintf("Invalid keerer")
	return sdk.ErrUnknownRequest(errMsg).Result()
}

func (m *orderModule) invalidKeeperTags() sdk.Tags {
	resTags := sdk.NewTags()
	resTags = resTags.AppendTag("order-endblocker", "Invalid keeper")
	return resTags
}