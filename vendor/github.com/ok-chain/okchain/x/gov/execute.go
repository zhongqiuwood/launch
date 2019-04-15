package gov

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Execute(ctx sdk.Context, gk Keeper, p Proposal) (err error) {
	switch p.GetProposalType() {
	case ProposalTypeParameterChange:
		return ParameterProposalExecute(ctx, gk, p.(*ParameterProposal))
	case ProposalTypeDexList:
		return DexListProposalExecute(ctx, gk, p.(*DexListProposal))
	}
	return nil
}

func ParameterProposalExecute(ctx sdk.Context, gk Keeper, pp *ParameterProposal) (err error) {
	ctx.Logger().Info("Execute ParameterProposal begin")
	curHeight := ctx.BlockHeight()
	if pp.Height > curHeight {
		gk.InsertWaitingProposalQueue(ctx, pp.ProposalID)
		return
	}

	for _, param := range pp.Params {
		paramSet, _ := gk.paramsKeeper.GetParamSet(param.Subspace)
		value, err := paramSet.ValidateKV(param.Key, param.Value)
		if err != nil {
			ctx.Logger().Error("Execute ParameterProposal Failed", "proposal", pp.ProposalID, "key", param.Key, "value", param.Value, "error", err)
		}

		subspace, found := gk.paramsKeeper.GetSubspace(param.Subspace)
		if found {
			fmt.Println(subspace.Name())
			subspace.Set(ctx, []byte(param.Key), value)
			ctx.Logger().Info("Execute ParameterProposal Successed", "proposal", pp.ProposalID, "key", param.Key, "value", param.Value, "VVV", value)
		} else {
			ctx.Logger().Error("Execute ParameterProposal Failed", "proposal", pp.ProposalID, "key", param.Key, "value", param.Value, "error", fmt.Sprintf("not found subspace(%s)", param.Subspace))
		}
	}
	gk.RemoveFromWaitingProposalQueue(ctx, pp.ProposalID)

	return
}


func DexListProposalExecute(ctx sdk.Context, gk Keeper, pp *DexListProposal) (err error) {
	ctx.Logger().Info("Execute ParameterProposal begin")
	curHeight := ctx.BlockHeight()
	if pp.BlockHeight > uint64(curHeight) {
		gk.InsertWaitingProposalQueue(ctx, pp.ProposalID)
		return
	}

	saveTokenPair(ctx, gk, pp)
	return
}