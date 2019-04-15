package backend

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type EndBlockEvent struct {
	ctx         sdk.Context
	blockHeight int64
	timestamp   int64
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

type Match struct {
	BlockHeight int64  `json:"blockHeight"`
	Product     string `gorm:"index;type:varchar(20)" json:"product"`
	Price       string `json:"price"`
	Quantity    string `json:"volume"`
	Timestamp   int64  `gorm:"index;type:int64" json:"timestamp"`
}

type Ticker struct {
	Product          string  `json:"product"`
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
	return fmt.Sprintf("[Ticker] Product: %s, TStr: %s, Timestamp: %d, OCHLV(%f, %f, %f, %f, %f) [%f, %f])",
		t.Product, timeString(t.Timestamp), t.Timestamp, t.Open, t.Close, t.High, t.Low, t.Volume, t.Change, t.ChangePercentage)
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
