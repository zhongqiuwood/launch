package backend

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ok-chain/okchain/x/order"
	"github.com/ok-chain/okchain/x/token"
)

// expected order keeper
type OrderKeeper interface {
	GetOrder(ctx sdk.Context, orderId string) *order.Order
	GetUpdatedOrderIds(ctx sdk.Context, blockHeight int64) []string
	GetBlockOrderNum(ctx sdk.Context, blockHeight int64) int64
	GetBlockMatchResult(ctx sdk.Context, blockHeight int64) *order.BlockMatchResult
}

// expected token keeper
type TokenKeeper interface {
	GetFeeDetailList(ctx sdk.Context, blockHeight int64) []token.FeeDetail
}
