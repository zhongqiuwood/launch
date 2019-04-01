package order

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

const (
	OrderStatusOpen      = 0
	OrderStatusFilled    = 1
	OrderStatusCancelled = 2
	OrderStatusExpired   = 3
)

type Order struct {
	OrderId        string         // order id, denoted as "blockHeight-orderNumInBlock".
	Sender         sdk.AccAddress // order maker address
	Product        string         // product for trading pair in full name of the tokens
	Side           string         // BUY/SELL
	Price          sdk.Dec        // price of the order, which is the real price multiplied by 1e8 (10^8) and rounded to integer
	Quantity       sdk.Dec        // quantity of the order, which is the real quantity multiplied by 1e8 (10^8) and rounded to integer
	Status         int64
	RemainQuantity sdk.Dec // Remaining quantity of the order, which is the real quantity multiplied by 1e8 (10^8) and rounded to integer
}

func FormatOrderId(blockHeight, orderNum int64) string {
	return fmt.Sprintf("%010d-%04d", blockHeight, orderNum)
}

func (order *Order) String() string {
	if orderJson, err := json.Marshal(order); err != nil {
		panic(err)
	} else {
		return string(orderJson)
	}
}

func (order *Order) Fill(amt sdk.Dec) {
	order.RemainQuantity = order.RemainQuantity.Sub(amt)
	if order.RemainQuantity.IsZero() {
		order.Status = OrderStatusFilled
	}
}

func (order *Order) FillAll() {
	order.RemainQuantity = sdk.ZeroDec()
	order.Status = OrderStatusFilled
}

func (order *Order) Cancel() {
	order.Status = OrderStatusCancelled
}

func (order *Order) Expire() {
	order.Status = OrderStatusExpired
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
