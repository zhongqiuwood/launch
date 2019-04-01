package token

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//
type MsgTokenIssue struct {
	Name   string
	Symbol string
	//OriginalSymbol string
	TotalSupply int64
	Owner       sdk.AccAddress
	Tokens      sdk.Coins
	//Mintable bool `json:"-"`
	Mintable bool
}

func NewMsgTokenIssue(name, symbol string, totalSupply int64, owner sdk.AccAddress, tokens sdk.Coins, mintable bool) MsgTokenIssue {
	return MsgTokenIssue{
		Name:   name,
		Symbol: symbol,
		//OriginalSymbol: originalSymbol,
		TotalSupply: totalSupply,
		Owner:       owner,
		Tokens:      tokens,
		Mintable:    mintable,
	}
}

func (msg MsgTokenIssue) Route() string { return "token" }

func (msg MsgTokenIssue) Type() string { return "issue" }

// ValidateBasic Implements Msg.
func (msg MsgTokenIssue) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Symbol) == 0 {
		return sdk.ErrUnknownRequest("Symbol cannot be empty")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTokenIssue) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgTokenIssue) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgTokenBurn burn coins
type MsgTokenBurn struct {
	Symbol string
	Amount int64
	Owner  sdk.AccAddress
}

func NewMsgTokenBurn(symbol string, amount int64, owner sdk.AccAddress) MsgTokenBurn {
	return MsgTokenBurn{
		Symbol: symbol,
		Amount: amount,
		Owner:  owner,
	}
}

func (msg MsgTokenBurn) Route() string { return "token" }

func (msg MsgTokenBurn) Type() string { return "burn" }

func (msg MsgTokenBurn) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Symbol) == 0 {
		return sdk.ErrUnknownRequest("Symbol cannot be empty")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTokenBurn) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgTokenBurn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgTokenBurn burn coins
type MsgTokenFreeze struct {
	Symbol string
	Amount int64
	Owner  sdk.AccAddress
}

func NewMsgTokenFreeze(symbol string, amount int64, owner sdk.AccAddress) MsgTokenFreeze {
	return MsgTokenFreeze{
		Symbol: symbol,
		Amount: amount,
		Owner:  owner,
	}
}

func (msg MsgTokenFreeze) Route() string { return "token" }

func (msg MsgTokenFreeze) Type() string { return "freeze" }

func (msg MsgTokenFreeze) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Symbol) == 0 {
		return sdk.ErrUnknownRequest("Symbol cannot be empty")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTokenFreeze) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgTokenFreeze) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgTokenBurn burn coins
type MsgTokenUnfreeze struct {
	Symbol string
	Amount int64
	Owner  sdk.AccAddress
}

func NewMsgTokenUnfreeze(symbol string, amount int64, owner sdk.AccAddress) MsgTokenUnfreeze {
	return MsgTokenUnfreeze{
		Symbol: symbol,
		Amount: amount,
		Owner:  owner,
	}
}

func (msg MsgTokenUnfreeze) Route() string { return "token" }

func (msg MsgTokenUnfreeze) Type() string { return "unfreeze" }

func (msg MsgTokenUnfreeze) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Symbol) == 0 {
		return sdk.ErrUnknownRequest("Symbol cannot be empty")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTokenUnfreeze) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgTokenUnfreeze) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

type MsgTokenMint struct {
	Symbol string
	Amount int64
	Owner  sdk.AccAddress
}

func NewMsgTokenMint(symbol string, amount int64, owner sdk.AccAddress) MsgTokenMint {
	return MsgTokenMint{
		Symbol: symbol,
		Amount: amount,
		Owner:  owner,
	}
}

func (msg MsgTokenMint) Route() string { return "token" }

func (msg MsgTokenMint) Type() string { return "mint" }

func (msg MsgTokenMint) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Symbol) == 0 {
		return sdk.ErrUnknownRequest("Symbol cannot be empty")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTokenMint) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgTokenMint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
