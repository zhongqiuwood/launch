package backend

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okchain/okdex/x/order"
)

// expected token keeper
type OrderKeeper interface {
	GetOrder(ctx sdk.Context, orderId string) *order.Order
	GetMatchResultMap(ctx sdk.Context, blockHeight int64) map[string]order.MatchResult
}
