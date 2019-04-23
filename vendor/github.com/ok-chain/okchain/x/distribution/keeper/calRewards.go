package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// calculate the total rewards accrued by a delegation
func (k Keeper) calDelegationRewards(ctx sdk.Context, valAddr sdk.ValAddress, delAddr sdk.AccAddress) (rewards sdk.DecCoins) {
	// a stake sanity check - recalculated final stake should be less than or equal to current stake
	// here we cannot use Equals because stake is truncated when multiplied by slash fractions
	// we could only use equals if we had arbitrary-precision rationals

	//outstanding := k.GetValidatorOutstandingRewards(ctx, valAddr)
	outstanding := k.GetValidatorCurrentRewards(ctx, valAddr).Rewards
	if outstanding.IsZero() {
		return sdk.NewDecCoins(sdk.Coins{})
	}

	del, f1 := k.GetDelegation(ctx, delAddr, valAddr)
	val, f2 := k.GetValidator(ctx, valAddr)
	//allshares := k.stakingKeeper.Validator(ctx, val.GetOperator()).GetDelegatorShares()

	if !f1 || !f2 {
		return sdk.NewDecCoins(sdk.Coins{})
	}
	//omit decimal
	rewards = outstanding.MulDecTruncate(del.GetShares()).QuoDecTruncate(val.GetDelegatorShares())
	return
}
