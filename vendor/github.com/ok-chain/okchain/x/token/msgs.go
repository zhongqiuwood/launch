package token

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//
type MsgTokenIssue struct {
	Name           string
	Symbol         string
	OriginalSymbol string
	TotalSupply    int64
	Owner          sdk.AccAddress
	//Tokens      sdk.Coins
	//Mintable bool `json:"-"`
	Mintable bool
}

func NewMsgTokenIssue(name, symbol, originalSymbol string, totalSupply int64, owner sdk.AccAddress, mintable bool) MsgTokenIssue {
	return MsgTokenIssue{
		Name:           name,
		Symbol:         symbol,
		OriginalSymbol: originalSymbol,
		TotalSupply:    totalSupply,
		Owner:          owner,
		//Tokens:      tokens,
		Mintable: mintable,
	}
}

func (msg MsgTokenIssue) Route() string { return "token" }

func (msg MsgTokenIssue) Type() string { return "issue" }

// ValidateBasic Implements Msg.
func (msg MsgTokenIssue) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.OriginalSymbol) == 0 {
		return sdk.ErrUnknownRequest("OriginalSymbol cannot be empty")
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
	Amount string
	Owner  sdk.AccAddress
}

func NewMsgTokenBurn(symbol string, amount string, owner sdk.AccAddress) MsgTokenBurn {
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
	Amount string
	Owner  sdk.AccAddress
}

func NewMsgTokenFreeze(symbol string, amount string, owner sdk.AccAddress) MsgTokenFreeze {
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
	Amount string
	Owner  sdk.AccAddress
}

func NewMsgTokenUnfreeze(symbol string, amount string, owner sdk.AccAddress) MsgTokenUnfreeze {
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

type MsgMultiSend struct {
	From      sdk.AccAddress
	Transfers []TransferUnit
}

func NewMsgMultiSend(from sdk.AccAddress, transfers []TransferUnit) MsgMultiSend {
	return MsgMultiSend{
		From:      from,
		Transfers: transfers,
	}
}

func (msg MsgMultiSend) Route() string { return "token" }

func (msg MsgMultiSend) Type() string { return "multi-send" }

func (msg MsgMultiSend) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress(msg.From.String())
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgMultiSend) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgMultiSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// MsgSend - high level transaction of the coin module
type MsgSend struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Amount      sdk.Coins      `json:"amount"`
}

func NewMsgTokenSend(from, to sdk.AccAddress, coins sdk.Coins) MsgSend {
	return MsgSend{
		FromAddress: from,
		ToAddress:   to,
		Amount:      coins,
	}
}

// Route Implements Msg.
func (msg MsgSend) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSend) Type() string { return "send" }

// ValidateBasic Implements Msg.
func (msg MsgSend) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress("missing recipient address")
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}
	if !msg.Amount.IsAllPositive() {
		return sdk.ErrInsufficientCoins("send amount must be positive")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSend) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
	//return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgSend - high level transaction of the coin module
type MsgTokenTransfer struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Symbol      string         `json:"symbol"`
}

func NewMsgTokenTransfer(from, to sdk.AccAddress, symbol string) MsgTokenTransfer {
	return MsgTokenTransfer{
		FromAddress: from,
		ToAddress:   to,
		Symbol:      symbol,
	}
}

// Route Implements Msg.
func (msg MsgTokenTransfer) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgTokenTransfer) Type() string { return "transfer" }

// ValidateBasic Implements Msg.
func (msg MsgTokenTransfer) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress("missing recipient address")
	}
	if msg.Symbol == "" {
		return sdk.ErrInvalidCoins("token is invalid: ")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTokenTransfer) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
	//return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgTokenTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}
