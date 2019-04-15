package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenPair struct {
	BaseAssetSymbol  string  `json:"base_asset_symbol"`
	QuoteAssetSymbol string  `json:"quote_asset_symbol"`
	InitPrice        sdk.Dec `json:"price"`
	MaxPriceDigit    int64   `json:"max_price_digit"`
	MaxQuantityDigit int64   `json:"max_size_digit"`
	//MergeTypes       string  `json:"merge_types"`
	MinQuantity sdk.Dec `json:"min_trade_size"`
}
