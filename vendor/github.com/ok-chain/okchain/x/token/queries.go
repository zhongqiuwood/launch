package token

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the governance Querier
const (
	QueryInfo       = "info"
	QueryTokens     = "tokens"
	QueryMarket     = "market"
	QueryParameters = "params"
	QueryCurrency   = "currency"
	QueryAccount    = "accounts"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryInfo:
			return queryInfo(ctx, path[1:], req, keeper)
		case QueryTokens:
			return queryTokens(ctx, path[1:], req, keeper)
		case QueryMarket:
			return queryMarket(ctx, path[1:], req, keeper)
		case QueryParameters:
			return queryParameters(ctx, keeper)
		case QueryCurrency:
			return queryCurrency(ctx, path[1:], req, keeper)
		case QueryAccount:
			return queryAccount(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown token query endpoint")
		}
	}
}

// nolint: unparam
func queryInfo(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	name := path[0]

	token := keeper.GetTokenInfo(ctx, name)

	if token.Symbol == "" {
		return nil, sdk.ErrInvalidCoins("unknown token")
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, token)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

func queryTokens(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	tokens := keeper.GetTokensInfo(ctx)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, tokens)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

func queryCurrency(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	tokens := keeper.GetCurrencysInfo(ctx)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, tokens)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

func queryMarket(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	tokenPairs := keeper.GetTokenPairs(ctx)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, tokenPairs)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

func queryAccount(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	addr, err1 := sdk.AccAddressFromBech32(path[0])
	if err1 != nil {
		return res, sdk.ErrInvalidAddress(path[0])
	}
	coinsInfo := keeper.GetCoinsInfo(ctx, addr)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, coinsInfo)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

func queryParameters(ctx sdk.Context, keeper Keeper) (res []byte, err sdk.Error) {
	params := keeper.GetParams(ctx)
	res, errRes := codec.MarshalJSONIndent(keeper.cdc, params)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}
