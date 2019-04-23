package keeper

import (
	`fmt`
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/ok-chain/okchain/x/distribution/types"
)

// allocate fees handles distribution of the collected fees

// allocate tokens to a particular validator, splitting according to commission
func (k Keeper) AllocateTokensToValidator(ctx sdk.Context, val sdk.Validator, tokens sdk.DecCoins) {

	// split tokens between validator and delegators according to commission
	commission := tokens.MulDec(val.GetCommission())
	shared := tokens.Sub(commission)

	// update current commission
	currentCommission := k.GetValidatorAccumulatedCommission(ctx, val.GetOperator())
	currentCommission = currentCommission.Add(commission)
	k.SetValidatorAccumulatedCommission(ctx, val.GetOperator(), currentCommission)

	// update current rewards
	currentRewards := k.GetValidatorCurrentRewards(ctx, val.GetOperator())
	currentRewards.Rewards = currentRewards.Rewards.Add(shared)
	k.SetValidatorCurrentRewards(ctx, val.GetOperator(), currentRewards)

	// update outstanding rewards
	outstanding := k.GetValidatorOutstandingRewards(ctx, val.GetOperator())
	outstanding = outstanding.Add(tokens)
	k.SetValidatorOutstandingRewards(ctx, val.GetOperator(), outstanding)
}

// allocate fees handles distribution of the collected fees, remaining to proposer
func (k Keeper) AllocateTokens(ctx sdk.Context, sumPrecommitPower, totalPower int64, proposer sdk.ConsAddress, votes []abci.VoteInfo) {

	// fetch collected fees & fee pool
	feesCollectedInt := k.feeCollectionKeeper.GetCollectedFees(ctx)
	feesCollected := sdk.NewDecCoins(feesCollectedInt)

	// clear collected fees, which will now be distributed
	k.feeCollectionKeeper.ClearCollectedFees(ctx)

	// pay proposer
	proposerValidator := k.stakingKeeper.ValidatorByConsAddr(ctx, proposer)
	if totalPower == 0 {
		k.AllocateTokensToValidator(ctx, proposerValidator, feesCollected)
		return
	}

	remaining := feesCollected
	// allocate tokens proportionally to voting power
	// TODO consider parallelizing later, ref https://github.com/cosmos/cosmos-sdk/pull/3099#discussion_r246276376
	for _, vote := range votes {
		validator := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)

		//validator not initialized handler, XXX, will be removed in the fulture
		if proposerValidator == nil{
			k.feeCollectionKeeper.AddCollectedFees(ctx, feesCollectedInt)
			return
		}

		// TODO consider microslashing for missing votes.
		// ref https://github.com/cosmos/cosmos-sdk/issues/2525#issuecomment-430838701
		powerFraction := sdk.NewDec(vote.Validator.Power).QuoTruncate(sdk.NewDec(totalPower))
		reward := feesCollected.MulDecTruncate(powerFraction)
		reward = reward.Intersect(remaining)
		k.AllocateTokensToValidator(ctx, validator, reward)
		remaining = remaining.Sub(reward)
	}
	// if has remaining, allocate to proposer
	k.AllocateTokensToValidator(ctx, proposerValidator, remaining)
}

func (k Keeper) distributeValRewards(ctx sdk.Context, valAddr sdk.ValAddress) {
	//all, commission, shared rewards
	outstanding := k.GetValidatorOutstandingRewards(ctx, valAddr)
	commission := k.GetValidatorAccumulatedCommission(ctx, valAddr)
	changes := sdk.DecCoins{{sdk.DefaultBondDenom, sdk.NewDec(0)}}
	if !commission.IsZero() {
		//substract from outstanding
		outstanding = outstanding.Sub(commission)

		//split into integral & remainder
		coins, remainder := commission.TruncateDecimal()
		changes.Add(remainder)

		//add to validator account
		if !coins.IsZero() {
			withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(valAddr))
			if _, _, err := k.bankKeeper.AddCoins(ctx, withdrawAddr, coins); err != nil {
				panic(err)
			}
		}
	}
	delegations := k.GetValidatorDelegations(ctx, valAddr) //Get delegations from snapshoot
	validator, _ := k.GetValidator(ctx, valAddr)
	allshares := validator.GetDelegatorShares()
	for _, del := range delegations {
		rewards := outstanding.MulDecTruncate(del.GetShares()).QuoDecTruncate(allshares)
		withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(del.GetDelegatorAddr()))
		coins, remainder := rewards.TruncateDecimal()

		//remainder to pool
		changes.Add(remainder)

		if _, _, err := k.bankKeeper.AddCoins(ctx, withdrawAddr, coins); err != nil {
			panic(err)
		}

		outstanding = outstanding.Sub(rewards)

	}

	if !outstanding.IsZero() {
		changes.Add(outstanding)
	}
	//adding changes to the validator, changes sum should be integral.
	coins, remainder := changes.TruncateDecimal()
	if !remainder.IsZero() {
		panic(fmt.Sprintf("lost some change of coins, losing okb: %+v", remainder))
	} else if !coins.IsZero() {
		withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(valAddr))
		if _, _, err := k.bankKeeper.AddCoins(ctx, withdrawAddr, coins); err != nil {
			panic(err)
		}
	}

	// set current rewards
	k.SetValidatorCurrentRewards(ctx, valAddr, types.NewValidatorCurrentRewards(sdk.DecCoins{}, 1))
	// set accumulated commission
	k.SetValidatorAccumulatedCommission(ctx, valAddr, types.InitialValidatorAccumulatedCommission())
	// set outstanding rewards
	k.SetValidatorOutstandingRewards(ctx, valAddr, sdk.DecCoins{})
}

//distribute rewards to validators and delegators
func (k Keeper) DistributeAllRewards(ctx sdk.Context) {
	//distribute all rewards, first snapshoot must be set in staking genesis, or this will panic!!!
	k.IterateValidators(ctx, func(_ int64, val sdk.Validator) (stop bool) {
		k.distributeValRewards(ctx, val.GetOperator())
		return false
	})
}
