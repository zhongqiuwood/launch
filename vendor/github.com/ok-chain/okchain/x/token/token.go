package token

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Token struct {
	Name           string         `json:"name"`            // token的名字
	Symbol         string         `json:"symbol"`          // token的唯一标识
	OriginalSymbol string         `json:"original_symbol"` // token的原始标识
	TotalSupply    int64          `json:"total_supply"`    // token的总量
	Owner          sdk.AccAddress `json:"owner"`           // token的所有者
	Mintable       bool           `json:"mintable"`        // token是否可以增发
}

func (token Token) String() string {
	b, _ := json.Marshal(token)
	return string(b)
}

type Currency struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	TotalSupply int64  `json:"total_supply"`
}

func (currency Currency) String() string {
	b, _ := json.Marshal(currency)
	return string(b)
}

type ByDenom sdk.Coins

func (d ByDenom) Len() int           { return len(d) }
func (d ByDenom) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByDenom) Less(i, j int) bool { return d[i].Denom < d[j].Denom }

type Transfer struct {
	To     string `json:"to"`
	Amount string `json:"amount"`
}

type TransferUnit struct {
	To    sdk.AccAddress `json:"to"`
	Coins sdk.Coins      `json:"coins"`
}

type CoinsInfo []CoinInfo

func (d CoinsInfo) Len() int           { return len(d) }
func (d CoinsInfo) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d CoinsInfo) Less(i, j int) bool { return d[i].Symbol < d[j].Symbol }

type CoinInfo struct {
	Symbol    string `json:"symbol"`
	Available string `json:"available"`
	Freeze    string `json:"freeze"`
	Locked    string `json:"locked"`
}
