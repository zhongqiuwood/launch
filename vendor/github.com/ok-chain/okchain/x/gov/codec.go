package gov

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var msgCdc = codec.New()

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSubmitProposal{}, "cosmos-sdk/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "cosmos-sdk/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "cosmos-sdk/MsgVote", nil)
	cdc.RegisterConcrete(MsgDexListSubmitProposal{}, "okdex/MsgDexListSubmitProposal", nil)
	cdc.RegisterConcrete(MsgDexList{}, "gov/DexList", nil)

	cdc.RegisterInterface((*Proposal)(nil), nil)
	cdc.RegisterConcrete(&TextProposal{}, "gov/TextProposal", nil)
	cdc.RegisterConcrete(&DexListProposal{}, "gov/DexListProposal", nil)
	cdc.RegisterConcrete(&ParameterProposal{}, "gov/ParameterProposal", nil)
}

func init() {
	RegisterCodec(msgCdc)
}
