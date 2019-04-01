package order

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type MsgNewOrder struct {
	Sender   sdk.AccAddress // order maker address
	Product  string         // product for trading pair in full name of the tokens
	Side     string         // BUY/SELL
	Price    string         // price of the order
	Quantity string         // quantity of the order
}

// NewMsgNewOrder is a constructor function for MsgNewOrder
func NewMsgNewOrder(sender sdk.AccAddress, product string, side string, price string, quantity string) MsgNewOrder {
	msgNewOrder := MsgNewOrder{
		Sender:   sender,
		Product:  product,
		Side:     side,
		Price:    price,
		Quantity: quantity,
	}

	return msgNewOrder
}

// Name Implements Msg.
func (msg MsgNewOrder) Route() string { return "order" }

// Type Implements Msg.
func (msg MsgNewOrder) Type() string { return "new" }

// ValdateBasic Implements Msg.
func (msg MsgNewOrder) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if len(msg.Product) == 0 {
		return sdk.ErrUnknownRequest("Product cannot be empty")
	}
	symbols := strings.Split(msg.Product, "_")
	if len(symbols) != 2 {
		return sdk.ErrUnknownRequest("Product in invalid format")
	}
	if symbols[0] == "okb" {
		return sdk.ErrUnknownRequest("Cannot use okb as base asset")
	}
	price, err1 := sdk.NewDecFromStr(msg.Price)
	quantity, err2 := sdk.NewDecFromStr(msg.Quantity)
	if err1 != nil || err2 != nil {
		return sdk.ErrUnknownRequest("Price/Quantity must be decimal string")
	}
	if price.IsNegative() || quantity.IsNegative() {
		return sdk.ErrUnknownRequest("Price/Quantity must be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgNewOrder) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgNewOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

type MsgCancelOrder struct {
	Sender  sdk.AccAddress
	OrderId string
}

// NewMsgCancelOrder is a constructor function for MsgCancelOrder
func NewMsgCancelOrder(sender sdk.AccAddress, orderId string) MsgCancelOrder {
	msgCancelOrder := MsgCancelOrder{
		Sender:  sender,
		OrderId: orderId,
	}
	return msgCancelOrder
}

// Name Implements Msg.
func (msg MsgCancelOrder) Route() string { return "order" }

// Type Implements Msg.
func (msg MsgCancelOrder) Type() string { return "cancel" }

// ValdateBasic Implements Msg.
func (msg MsgCancelOrder) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCancelOrder) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgCancelOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
