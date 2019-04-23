package order

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

const (
	OrderStatusOpen                   = 0
	OrderStatusFilled                 = 1
	OrderStatusCancelled              = 2
	OrderStatusExpired                = 3
	OrderStatusPartialFilledCancelled = 4
	OrderStatusPartialFilledExpired   = 5
)

const (
	OrderExtraInfoKeyNewFee    = "newFee"
	OrderExtraInfoKeyCancelFee = "cancelFee"
	OrderExtraInfoKeyExpireFee = "expireFee"
	OrderExtraInfoKeyDealFee   = "dealFee"
)

type Order struct {
	TxHash         string         `json:"txHash"`         // txHash of the place order tx
	OrderId        string         `json:"orderId"`        // order id, denoted as "blockHeight-orderNumInBlock".
	Sender         sdk.AccAddress `json:"sender"`         // order maker address
	Product        string         `json:"product"`        // product for trading pair in full name of the tokens
	Side           string         `json:"side"`           // BUY/SELL
	Price          sdk.Dec        `json:"price"`          // price of the order
	Quantity       sdk.Dec        `json:"quantity"`       // quantity of the order
	Status         int64          `json:"status"`         // order status, (0-5) respectively represents (Open, Filled, Cancelled, Expired, PartialFilledCancelled, PartialFilledExpired)
	FilledAvgPrice sdk.Dec        `json:"filledAvgPrice"` // filled average price
	RemainQuantity sdk.Dec        `json:"remainQuantity"` // Remaining quantity of the order
	Timestamp      int64          `json:"timestamp"`      // created timestamp
	ExtraInfo      string         `json:"extraInfo"`      // extra info of order, json format, eg.{"cancelFee": "0.002okb"}
}

func (order *Order) String() string {
	if orderJson, err := json.Marshal(order); err != nil {
		panic(err)
	} else {
		return string(orderJson)
	}
}

func (order *Order) SetExtraInfo(extra map[string]string) {
	if extra != nil {
		bz, _ := json.Marshal(extra)
		order.ExtraInfo = string(bz)
	}
}

func (order *Order) GetExtraInfo() map[string]string {
	extra := make(map[string]string)
	if order.ExtraInfo != "" && order.ExtraInfo != "{}" {
		json.Unmarshal([]byte(order.ExtraInfo), &extra)
	}
	return extra
}

func (order *Order) SetExtraInfoWithKeyValue(key, value string) {
	extra := order.GetExtraInfo()
	extra[key] = value
	order.SetExtraInfo(extra)
}

func (order *Order) GetExtraInfoWithKey(key string) string {
	extra := order.GetExtraInfo()
	if value, ok := extra[key]; ok {
		return value
	}
	return ""
}

func (order *Order) RecordOrderNewFee(fee sdk.DecCoins) {
	order.SetExtraInfoWithKeyValue(OrderExtraInfoKeyNewFee, fee.String())
}

func (order *Order) RecordOrderCancelFee(fee sdk.DecCoins) {
	order.SetExtraInfoWithKeyValue(OrderExtraInfoKeyCancelFee, fee.String())
}

func (order *Order) RecordOrderExpireFee(fee sdk.DecCoins) {
	order.SetExtraInfoWithKeyValue(OrderExtraInfoKeyExpireFee, fee.String())
}

// An order may have several deals
func (order *Order) RecordOrderDealFee(fee sdk.DecCoins) {
	oldValue := order.GetExtraInfoWithKey(OrderExtraInfoKeyDealFee)
	if oldValue == "" {
		order.SetExtraInfoWithKeyValue(OrderExtraInfoKeyDealFee, fee.String())
		return
	}
	oldFee, _ := sdk.ParseDecCoins(oldValue)
	newFee := oldFee.Add(fee)
	order.SetExtraInfoWithKeyValue(OrderExtraInfoKeyDealFee, newFee.String())
}

func (order *Order) Fill(price, amt sdk.Dec) {
	filledSum := order.FilledAvgPrice.Mul(order.Quantity.Sub(order.RemainQuantity))
	newFilledSum := filledSum.Add(price.Mul(amt))
	order.RemainQuantity = order.RemainQuantity.Sub(amt)
	order.FilledAvgPrice = newFilledSum.Quo(order.Quantity.Sub(order.RemainQuantity))
	if order.RemainQuantity.IsZero() {
		order.Status = OrderStatusFilled
	}
}

func (order *Order) Cancel() {
	if order.RemainQuantity.Equal(order.Quantity) {
		order.Status = OrderStatusCancelled
	} else {
		order.Status = OrderStatusPartialFilledCancelled
	}
}

func (order *Order) Expire() {
	if order.RemainQuantity.Equal(order.Quantity) {
		order.Status = OrderStatusExpired
	} else {
		order.Status = OrderStatusPartialFilledExpired
	}
}

// when place a new order, we should lock the coins of sender
func (order *Order) NeedLockCoins() sdk.DecCoins {
	if order.Side == "BUY" {
		token := strings.Split(order.Product, "_")[1]
		amount := order.Price.Mul(order.Quantity)
		return sdk.DecCoins{{token, amount}}
	} else {
		token := strings.Split(order.Product, "_")[0]
		amount := order.Quantity
		return sdk.DecCoins{{token, amount}}
	}
}

// when order be cancelled/expired, we should unlock the coins of sender
func (order *Order) NeedUnlockCoins() sdk.DecCoins {
	if order.Side == "BUY" {
		token := strings.Split(order.Product, "_")[1]
		amount := order.Price.Mul(order.RemainQuantity)
		return sdk.DecCoins{{token, amount}}
	} else {
		token := strings.Split(order.Product, "_")[0]
		amount := order.RemainQuantity
		return sdk.DecCoins{{token, amount}}
	}
}
