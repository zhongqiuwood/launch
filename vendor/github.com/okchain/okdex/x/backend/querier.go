package backend

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the governance Querier
const (
	QueryTradeHistory = "trades"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryTradeHistory:
			return queryTradeHistory(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown backend endpoint")
		}
	}
}

// nolint: unparam
func queryTradeHistory(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	trades := keeper.GetTrades(ctx, path[0])
	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, trades)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}
