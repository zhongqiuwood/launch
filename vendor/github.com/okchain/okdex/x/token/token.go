package token

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Token struct {
	Name           string         `json:"name"`
	Symbol         string         `json:"symbol"`
	OriginalSymbol string         `json:"original_symbol"`
	TotalSupply    int64          `json:"total_supply"`
	Owner          sdk.AccAddress `json:"owner"`
	//Mintable       bool           `json:"-"`
	Mintable bool `json:"mintable"`
}

func (token Token) String() string {
	b, _ := json.Marshal(token)
	return string(b)
}

type ByDenom sdk.Coins

func (d ByDenom) Len() int           { return len(d) }
func (d ByDenom) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByDenom) Less(i, j int) bool { return d[i].Denom < d[j].Denom }
