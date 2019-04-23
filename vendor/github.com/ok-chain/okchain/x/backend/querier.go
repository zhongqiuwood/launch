package backend

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ok-chain/okchain/x/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"sort"
)

// query endpoints supported by the governance Querier
const (
	QueryDealList   = "deals"
	QueryFeeDetails = "fees"
	QueryOrderList  = "orders"
	QueryTxList     = "txs"
	QueryCandleList = "candles"
	QueryTickerList = "tickers"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 50
)

func GetPage(page, perPage int) (offset, limit int) {
	o := (page - 1) * perPage
	l := perPage
	return o, l
}

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryDealList:
			return queryDeals(ctx, path[1:], req, keeper)
		case QueryFeeDetails:
			return queryFeeDetails(ctx, path[1:], req, keeper)
		case QueryOrderList:
			return queryOrderList(ctx, path[1:], req, keeper)
		case QueryTxList:
			return queryTxList(ctx, path[1:], req, keeper)
		case QueryCandleList:
			return queryCandleList(ctx, path[1:], req, keeper)
		case QueryTickerList:
			return queryTickerList(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown backend endpoint")
		}
	}
}

type QueryDealsParams struct {
	Address string
	Product string
	Page    int
	PerPage int
}

func NewQueryDealsParams(addr, product string, page, perPage int) QueryDealsParams {
	if page == 0 && perPage == 0 {
		page = DefaultPage
		perPage = DefaultPerPage
	}
	return QueryDealsParams{
		Address: addr,
		Product: product,
		Page:    page,
		PerPage: perPage,
	}
}

// nolint: unparam
func queryDeals(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryDealsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}
	offset, limit := common.GetPage(params.Page, params.PerPage)
	deals, total := keeper.GetDeals(ctx, params.Address, params.Product, offset, limit)
	response := common.GetListResponse(total, params.Page, params.PerPage, *deals)
	bz, err := json.Marshal(response)
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

type QueryFeeDetailsParams struct {
	Address string
	Page    int
	PerPage int
}

// creates a new instance of NewQueryOrderListParams
func NewQueryFeeDetailsParams(addr string, page, perPage int) QueryFeeDetailsParams {
	if page == 0 && perPage == 0 {
		page = DefaultPage
		perPage = DefaultPerPage
	}
	return QueryFeeDetailsParams{
		Address: addr,
		Page:    page,
		PerPage: perPage,
	}
}

// nolint: unparam
func queryFeeDetails(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryFeeDetailsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	ctx.Logger().Debug(fmt.Sprintf("queryFeeDetails params : %+v", params))

	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}
	offset, limit := GetPage(params.Page, params.PerPage)
	feeDetails, total := keeper.GetFeeDetails(ctx, path[0], offset, limit)
	response := common.GetListResponse(total, params.Page, params.PerPage, *feeDetails)
	bz, err := json.Marshal(response)
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

type QueryKlinesParams struct {
	Product     string
	Granularity int
	Size        int
}

func NewQueryKlinesParams(product string, granularity, size int) QueryKlinesParams {
	return QueryKlinesParams{
		product,
		granularity,
		size,
	}
}

func queryCandleList(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {

	var params QueryKlinesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}
	ctx.Logger().Debug(fmt.Sprintf("queryCandleList : %+v", params))
	restData, err := keeper.GetCandles(params.Product, params.Granularity, params.Size)
	finalResult := map[string]interface{}{
		"code":      0,
		"data":      nil,
		"detailMsg": "",
		"msg":       "",
	}
	if err != nil {
		finalResult["detailMsg"] = err.Error()
	} else {
		finalResult["data"] = restData
	}

	bz, err := json.Marshal(finalResult)

	return bz, nil
}

type QueryTickerParams struct {
	Product string `json:"product"`
	Count   int    `json:"count"`
	Sort    bool   `json:"sort"`
}

func queryTickerList(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	params := QueryTickerParams{}
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data, ", err.Error()))
	}

	products := []string{}
	if params.Product != "" {
		products = append(products, params.Product)
	}

	tickers := keeper.GetTickers(products, params.Count)

	var sortedTickers Tickers = tickers
	sort.Sort(sortedTickers)

	finalResult := map[string]interface{}{
		"code":      0,
		"data":      sortedTickers,
		"detailMsg": "",
		"msg":       "",
	}

	bz, _ := json.Marshal(finalResult)
	return bz, nil

}

type QueryOrderListParams struct {
	Address string
	Product string
	Page    int
	PerPage int
}

// creates a new instance of NewQueryOrderListParams
func NewQueryOrderListParams(addr, product string, page, perPage int) QueryOrderListParams {
	if page == 0 && perPage == 0 {
		page = DefaultPage
		perPage = DefaultPerPage
	}
	return QueryOrderListParams{
		Address: addr,
		Product: product,
		Page:    page,
		PerPage: perPage,
	}
}

func queryOrderList(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	isOpen := path[0] == "open"
	var params QueryOrderListParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}
	_, err = sdk.AccAddressFromBech32(params.Address)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("invalid address", err.Error()))
	}
	offset, limit := GetPage(params.Page, params.PerPage)
	orders, total := keeper.GetOrderList(ctx, params.Address, params.Product, isOpen, offset, limit)

	response := common.GetListResponse(total, params.Page, params.PerPage, *orders)
	bz, err := json.Marshal(response)
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

type QueryTxListParams struct {
	Address   string
	TxType    int64
	StartTime int64
	EndTime   int64
	Page      int
	PerPage   int
}

// creates a new instance of NewQueryOrderListParams
func NewQueryTxListParams(addr string, txType, startTime, endTime int64, page, perPage int) QueryTxListParams {
	if page == 0 && perPage == 0 {
		page = DefaultPage
		perPage = DefaultPerPage
	}
	return QueryTxListParams{
		Address:   addr,
		TxType:    txType,
		StartTime: startTime,
		EndTime:   endTime,
		Page:      page,
		PerPage:   perPage,
	}
}

func queryTxList(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryTxListParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}
	_, err = sdk.AccAddressFromBech32(params.Address)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("invalid address", err.Error()))
	}
	offset, limit := GetPage(params.Page, params.PerPage)
	txs, total := keeper.GetTransactionList(ctx, params.Address, params.TxType, params.StartTime, params.EndTime, offset, limit)

	response := common.GetListResponse(total, params.Page, params.PerPage, *txs)
	bz, err := json.Marshal(response)
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}
