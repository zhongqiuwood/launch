package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenPair struct {
	BaseAssetSymbol  string  `json:"base_asset_symbol"`		// 基础货币
	QuoteAssetSymbol string  `json:"quote_asset_symbol"`	// 报价货币
	InitPrice        sdk.Dec `json:"price"`					// 价格
	MaxPriceDigit    int64   `json:"max_price_digit"`	 	// 最大交易价格的小数点位数
	MaxQuantityDigit int64   `json:"max_size_digit"`		// 最大交易数量的小数点位数
	MinQuantity      sdk.Dec `json:"min_trade_size"`		// 最小交易数量
}
