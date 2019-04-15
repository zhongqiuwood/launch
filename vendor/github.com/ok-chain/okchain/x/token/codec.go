package token

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgTokenIssue{}, "token/TokenIssue", nil)
	cdc.RegisterConcrete(MsgTokenBurn{}, "token/TokenBurn", nil)
	cdc.RegisterConcrete(MsgTokenFreeze{}, "token/TokenFreeze", nil)
	cdc.RegisterConcrete(MsgTokenUnfreeze{}, "token/TokenUnfreeze", nil)
	cdc.RegisterConcrete(MsgTokenMint{}, "token/TokenMint", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "token/MultiSend", nil)
	cdc.RegisterConcrete(MsgSend{}, "token/Send", nil)
	cdc.RegisterConcrete(MsgTokenTransfer{}, "token/Transfer", nil)
}
