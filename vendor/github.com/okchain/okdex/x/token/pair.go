package token

type TokenPair struct {
	BaseAssetSymbol  string `json:"base_asset_symbol"`
	QuoteAssetSymbol string `json:"quote_asset_symbol"`
	Price            string `json:"price"`
	TickSize         string `json:"tick_size"`
	LotSize          string `json:"lot_size"`
	MaxPriceDigit    uint64 `json:"max_price_digit"`
	MaxSizeDigit     uint64 `json:"max_size_digit"`
	MergeTypes       string `json:"merge_types"`
	MinTradeSize     string `json:"min_trade_size"`
}
