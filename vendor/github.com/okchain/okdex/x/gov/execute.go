package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Execute(ctx sdk.Context, gk Keeper, p Proposal) (err error) {
	switch p.GetProposalType() {
	case ProposalTypeParameterChange:
		return ParameterProposalExecute(ctx, gk, p.(*ParameterProposal))
	}
	return nil
}

func ParameterProposalExecute(ctx sdk.Context, gk Keeper, pp *ParameterProposal) (err error) {
	ctx.Logger().Info("Execute ParameterProposal begin")
	for _, param := range pp.Params {
		paramSet, _ := gk.paramsKeeper.GetParamSet(param.Subspace)
		value, err := paramSet.ValidateKV(param.Key, param.Value)
		if err != nil {
			ctx.Logger().Error("Execute ParameterProposal Failed", "key", param.Key, "value", param.Value, "error", err)
		}

		subspace, found := gk.paramsKeeper.GetSubspace(param.Subspace)
		if found {
			subspace.Set(ctx, []byte(param.Key), value)
			ctx.Logger().Info("Execute ParameterProposal Successed", "key", param.Key, "value", param.Value)
		} else {
			ctx.Logger().Info("Execute ParameterProposal Failed", "key", param.Key, "value", param.Value)
		}

	}

	return
}
