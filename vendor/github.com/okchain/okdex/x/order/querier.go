package order

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the governance Querier
const (
	QueryOrderDetail    = "detail"
	QueryDepthBook      = "depthbook"
	QueryMatchResultMap = "match"
	QueryParameters     = "params"
)

const DefaultBookSize = 200

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryOrderDetail:
			return queryOrder(ctx, path[1:], req, keeper)
		case QueryDepthBook:
			return queryDepthBook(ctx, path[1:], req, keeper)
		case QueryMatchResultMap:
			return queryMatchResultMap(ctx, path[1:], req, keeper)
		case QueryParameters:
			return queryParameters(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown order query endpoint")
		}
	}
}

// nolint: unparam
func queryOrder(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	order := keeper.GetOrder(ctx, path[0])
	if order == nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("order(%v) does not exist", path[0]))
	}
	bz, err2 := keeper.cdc.MarshalJSON(order)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

type QueryDepthBookParams struct {
	Product string
	Size    int
}

// creates a new instance of QueryProposalParams
func NewQueryDepthBookParams(product string, size int) QueryDepthBookParams {
	if size == 0 {
		size = DefaultBookSize
	}
	return QueryDepthBookParams{
		Product: product,
		Size:    size,
	}
}

type BookResItem struct {
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

type BookRes struct {
	Code int64         `json:"code"`
	Msg  string        `json:"msg"`
	Asks []BookResItem `json:"asks"`
	Bids []BookResItem `json:"bids"`
}

func (book *BookRes) String() string {
	if bookJson, err := json.Marshal(book); err != nil {
		panic(err)
	} else {
		return string(bookJson)
	}
}

// nolint: unparam
func queryDepthBook(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryDepthBookParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}
	depthBook := keeper.GetDepthBook(ctx, params.Product)

	var asks []BookResItem
	var bids []BookResItem
	for _, item := range *depthBook {
		if item.SellQuantity.IsPositive() {
			asks = append([]BookResItem{{item.Price.String(), item.SellQuantity.String()}}, asks...)
		}
		if item.BuyQuantity.IsPositive() {
			bids = append(bids, BookResItem{item.Price.String(), item.BuyQuantity.String()})
		}
	}
	if len(asks) > params.Size {
		asks = asks[:params.Size]
	}
	if len(bids) > params.Size {
		bids = bids[:params.Size]
	}

	bookRes := BookRes{
		Code: 0,
		Msg:  "",
		Asks: asks,
		Bids: bids,
	}
	bz, err2 := keeper.cdc.MarshalJSON(bookRes)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

// nolint: unparam
func queryMatchResultMap(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	blockHeight, _ := strconv.ParseInt(path[0], 10, 64)
	matchResultMap := keeper.GetMatchResultMap(ctx, blockHeight)
	bz, err2 := keeper.cdc.MarshalJSON(matchResultMap)
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
