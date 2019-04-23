package backend

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

const (
	TxTypeTransfer    = 1
	TxTypeOrderNew    = 2
	TxTypeOrderCancel = 3
)

const (
	TxSideBuy  = 1
	TxSideSell = 2
	TxSideFrom = 3
	TxSideTo   = 4
)

type EndBlockEvent struct {
	ctx         sdk.Context
	blockHeight int64
	timestamp   int64
}

type TxEvent struct {
	ctx       sdk.Context
	tx        *auth.StdTx
	txHash    string
	timestamp int64
}

type Deal struct {
	Timestamp   int64   `gorm:"index;type:int64" json:"timestamp"`
	BlockHeight int64   `gorm:"PRIMARY_KEY;type:int64" json:"blockHeight"`
	OrderId     string  `gorm:"PRIMARY_KEY;type:varchar(30)" json:"orderId"`
	Sender      string  `gorm:"index;type:varchar(80)" json:"sender"`
	Product     string  `gorm:"index;type:varchar(20)" json:"product"`
	Side        string  `gorm:"type:varchar(10)" json:"side"`
	Price       float64 `gorm:"type:DOUBLE" json:"price"`
	Quantity    float64 `gorm:"type:DOUBLE" json:"volume"`
}

type Ticker struct {
	Symbol           string  `json:"symbol"`
	CurrencyId       string  `json:"currency_id"`
	Timestamp        int64   `json:"timestamp"`
	Open             float64 `json:"open"`              // Open In 24h
	Close            float64 `json:"close"`             // Close in 24h
	High             float64 `json:"high"`              // High in 24h
	Low              float64 `json:"low"`               // Low in 24h
	Volume           float64 `json:"volume"`            // Volume in 24h
	Change           float64 `json:"change"`            // (Close - Open)
	ChangePercentage float64 `json:"change_percentage"` // Change / Open * 100%
}

func (t *Ticker) PrettyString() string {
	return fmt.Sprintf("[Ticker] Symbol: %s, TStr: %s, Timestamp: %d, OCHLV(%f, %f, %f, %f, %f) [%f, %f])",
		t.Symbol, timeString(t.Timestamp), t.Timestamp, t.Open, t.Close, t.High, t.Low, t.Volume, t.Change, t.ChangePercentage)
}

type Tickers []Ticker

func (tickers Tickers) Len() int {
	return len(tickers)
}

func (c Tickers) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (tickers Tickers) Less(i, j int) bool {
	return tickers[i].Change < tickers[j].Change
}

type KlineSnapShot struct {
	LastBlockHeight int64 `json:"last_block_height"`
	LastSyncedTime  int64 `json:"last_synced_time"`
}

type Order struct {
	TxHash         string `gorm:"type:varchar(80)" json:"txHash"`
	OrderId        string `gorm:"PRIMARY_KEY;type:varchar(30)" json:"orderId"`
	Sender         string `gorm:"index;type:varchar(80)" json:"sender"`
	Product        string `gorm:"index;type:varchar(20)" json:"product"`
	Side           string `gorm:"type:varchar(10)" json:"side"`
	Price          string `gorm:"type:varchar(40)" json:"price"`
	Quantity       string `gorm:"type:varchar(40)" json:"quantity"`
	Status         int64  `gorm:"index;type:int64" json:"status"`
	FilledAvgPrice string `gorm:"type:varchar(40)" json:"filledAvgPrice"`
	RemainQuantity string `gorm:"type:varchar(40)" json:"remainQuantity"`
	Timestamp      int64  `gorm:"index;type:int64" json:"timestamp"`
}

type Transaction struct {
	TxHash    string `gorm:"type:varchar(80)" json:"txHash"`
	Type      int64  `gorm:"index;type:int64" json:"type"` // 1:Transfer, 2:NewOrder, 3:CancelOrder
	Address   string `gorm:"index;type:varchar(80)" json:"address"`
	Symbol    string `gorm:"type:varchar(20)" json:"symbol"`
	Side      int64  `gorm:"type:int64" json:"side"` // 1:buy, 2:sell, 3:from, 4:to
	Quantity  string `gorm:"type:varchar(40)" json:"quantity"`
	Fee       string `gorm:"type:varchar(40)" json:"fee"`
	Timestamp int64  `gorm:"index;type:int64" json:"timestamp"`
}
