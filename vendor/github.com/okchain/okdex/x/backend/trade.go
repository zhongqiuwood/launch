package backend

import "github.com/jinzhu/gorm"

type Trade struct {
	BlockHeight int64  `json:"blockHeight"`
	OrderId     string `gorm:"type:varchar(30)" json:"orderId"`
	Sender      string `gorm:"index;type:varchar(80)" json:"sender"`
	Product     string `gorm:"index;type:varchar(20)" json:"product"`
	Side        string `gorm:"index;type:varchar(10)" json:"side"`
	Price       string `json:"price"`
	Quantity    string `json:"volume"`
	Timestamp   int64  `gorm:"index" json:"timestamp"`
}

type Match struct {
	BlockHeight int64  `json:"blockHeight"`
	Product     string `gorm:"index;type:varchar(20)" json:"product"`
	Price       string `json:"price"`
	Quantity    string `json:"volume"`
	Timestamp   int64  `gorm:"index" json:"timestamp"`
}

type KLineMin struct {
	*gorm.Model
	Product   string `gorm:"index;type:varchar(20)" json:"product"`
	Open      string `json:"open"`
	Close     string `json:"open"`
	High      string `json:"open"`
	Low       string `json:"open"`
	Volume    string `json:"open"`
	Timestamp int64  `gorm:"index" json:"timestamp"`
}
