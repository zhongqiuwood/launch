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
}

func (order *Order) String() string {
	if orderJson, err := json.Marshal(order); err != nil {
		panic(err)
	} else {
		return string(orderJson)
	}
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
