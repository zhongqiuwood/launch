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

func ParameterProposalExecute(ctx sdk.Context, gk Keeper, pp *ParameterProposal) (err sdk.Error) {
	ctx.Logger().Info("Execute ParameterProposal begin")
	curHeight := ctx.BlockHeight()
	if pp.Height > curHeight {
		gk.InsertWaitingProposalQueue(ctx, uint64(pp.Height), pp.ProposalID)
		return nil
	}

	for _, param := range pp.Params {
		subspace, found := gk.paramsKeeper.GetSubspace(param.Subspace)
		if found {
			paramSet, _ := gk.paramsKeeper.GetParamSet(param.Subspace)
			value, err := paramSet.ValidateKV(param.Key, param.Value)
			if err != nil {
				ctx.Logger().Error("Execute ParameterProposal Failed", "proposal", pp.ProposalID, "key", param.Key, "value", param.Value, "error", err)
				return err
			}
			subspace.Set(ctx, []byte(param.Key), value)
			ctx.Logger().Info("Execute ParameterProposal Successed", "proposal", pp.ProposalID, "key", param.Key, "value", param.Value, "VVV", value)
		} else {
			ctx.Logger().Error("Execute ParameterProposal Failed", "proposal", pp.ProposalID, "key", param.Key, "value", param.Value, "error", fmt.Sprintf("not found subspace(%s)", param.Subspace))
			return sdk.NewError(DefaultParamspace, CodeInvalidParamSubspace, "Param subspace %s not found", param.Subspace)
		}
	}
	gk.RemoveFromWaitingProposalQueue(ctx, uint64(pp.Height), pp.ProposalID)

	return nil
}

func DexListProposalExecute(ctx sdk.Context, gk Keeper, pp *DexListProposal) (err error) {
	ctx.Logger().Info("Execute DexListProposal begin")
	curHeight := ctx.BlockHeight()
	if pp.BlockHeight > uint64(curHeight) {
		gk.InsertWaitingProposalQueue(ctx, pp.BlockHeight, pp.ProposalID)
		return
	}

	res := saveTokenPair(ctx, gk, pp)
	if !res.IsOK() {
		// TODO:should not happen panic
		return fmt.Errorf("DexListProposalExecute saveTokenPair failed")
	}
	gk.RemoveFromWaitingProposalQueue(ctx, pp.BlockHeight, pp.ProposalID)
	return
}
