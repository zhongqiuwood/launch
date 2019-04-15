package keeper

import (
	"fmt"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ok-chain/okchain/x/distribution/types"
)

// allocate fees handles distribution of the collected fees
func (k Keeper) AllocateTokens(ctx sdk.Context, sumPrecommitPower, totalPower int64, proposer sdk.ConsAddress, votes []abci.VoteInfo) {
	logger := ctx.Logger().With("module", "x/distribution")

	// fetch collected fees & fee pool
	feesCollectedInt := k.feeCollectionKeeper.GetCollectedFees(ctx)
	feesCollected := sdk.NewDecCoins(feesCollectedInt)
	feePool := k.GetFeePool(ctx)

	// clear collected fees, which will now be distributed
	k.feeCollectionKeeper.ClearCollectedFees(ctx)

	// temporary workaround to keep CanWithdrawInvariant happy
	// general discussions here: https://github.com/cosmos/cosmos-sdk/issues/2906#issuecomment-441867634
	if totalPower == 0 {
		feePool.CommunityPool = feePool.CommunityPool.Add(feesCollected)
		k.SetFeePool(ctx, feePool)
		return
	}

	// calculate fraction votes
	fractionVotes := sdk.NewDec(sumPrecommitPower).Quo(sdk.NewDec(totalPower))

	// calculate proposer reward
	baseProposerReward := k.GetBaseProposerReward(ctx)
	bonusProposerReward := k.GetBonusProposerReward(ctx)
	proposerMultiplier := baseProposerReward.Add(bonusProposerReward.MulTruncate(fractionVotes))
	proposerReward := feesCollected.MulDecTruncate(proposerMultiplier)

	// pay proposer
	remaining := feesCollected
	proposerValidator := k.stakingKeeper.ValidatorByConsAddr(ctx, proposer)
	if proposerValidator != nil {
		k.AllocateTokensToValidator(ctx, proposerValidator, proposerReward)
		remaining = remaining.Sub(proposerReward)
	} else {
		// proposer can be unknown if say, the unbonding period is 1 block, so
		// e.g. a validator undelegates at block X, it's removed entirely by
		// block X+1's endblock, then X+2 we need to refer to the previous
		// proposer for X+1, but we've forgotten about them.
		logger.Error(fmt.Sprintf(
			"WARNING: Attempt to allocate proposer rewards to unknown proposer %s. "+
				"This should happen only if the proposer unbonded completely within a single block, "+
				"which generally should not happen except in exceptional circumstances (or fuzz testing). "+
				"We recommend you investigate immediately.",
			proposer.String()))
	}

	// calculate fraction allocated to validators
	communityTax := k.GetCommunityTax(ctx)
	voteMultiplier := sdk.OneDec().Sub(proposerMultiplier).Sub(communityTax)

	// allocate tokens proportionally to voting power
	// TODO consider parallelizing later, ref https://github.com/cosmos/cosmos-sdk/pull/3099#discussion_r246276376
	for _, vote := range votes {
		validator := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)

		// TODO consider microslashing for missing votes.
		// ref https://github.com/cosmos/cosmos-sdk/issues/2525#issuecomment-430838701
		powerFraction := sdk.NewDec(vote.Validator.Power).QuoTruncate(sdk.NewDec(totalPower))
		reward := feesCollected.MulDecTruncate(voteMultiplier).MulDecTruncate(powerFraction)
		reward = reward.Intersect(remaining)
		k.AllocateTokensToValidator(ctx, validator, reward)
		remaining = remaining.Sub(reward)
	}

	// allocate community funding
	feePool.CommunityPool = feePool.CommunityPool.Add(remaining)
	k.SetFeePool(ctx, feePool)
}

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

//distribute rewards to validators and delegators
func (k Keeper) DistributeAllRewards(ctx sdk.Context) {
	//distribute all rewards, first snapshoot must be set in staking genesis, or this will panic!!!
	k.IterateValidators(ctx, func(_ int64, val sdk.Validator) (stop bool) {
		k.distributeValRewards(ctx, val.GetOperator())
		return false
	})
}

func (k Keeper) distributeValRewards(ctx sdk.Context, valAddr sdk.ValAddress) {
	//all, commission, shared rewards
	outstanding := k.GetValidatorOutstandingRewards(ctx, valAddr)
	commission := k.GetValidatorAccumulatedCommission(ctx, valAddr)
	//curRewards := k.GetValidatorCurrentRewards(ctx, valAddr) //rewards excluding commssion
	if !commission.IsZero() {
		//substract from outstanding
		outstanding = outstanding.Sub(commission)

		//split into integral & remainder
		coins, remainder := commission.TruncateDecimal()

		//remainder to community pool
		feelPool := k.GetFeePool(ctx)
		feelPool.CommunityPool = feelPool.CommunityPool.Add(remainder)
		k.SetFeePool(ctx, feelPool)

		//add to validator account
		if !coins.IsZero() {
			withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(valAddr))
			if _, _, err := k.bankKeeper.AddCoins(ctx, withdrawAddr, coins); err != nil {
				panic(err)
			}
		}
		fmt.Printf("distribute commission[%+v] to validator[%+v], remainder[%+v] to communityPool\n", commission, valAddr.String(), remainder)
	}
	delegations := k.GetValidatorDelegations(ctx, valAddr) //Get delegations from snapshoot
	//delegations := keeper.stakingKeeper.GetValidatorDelegations(ctx,valAddr) //Get latest delegations
	//allshares := k.stakingKeeper.Validator(ctx, valAddr).GetDelegatorShares()
	validator, _ := k.GetValidator(ctx, valAddr)
	allshares := validator.GetDelegatorShares()
	for _, del := range delegations {
		rewards := outstanding.MulDecTruncate(del.GetShares()).QuoDecTruncate(allshares)
		withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(del.GetDelegatorAddr()))
		coins, remainder := rewards.TruncateDecimal()

		//remainder to community pool
		feelPool := k.GetFeePool(ctx)
		feelPool.CommunityPool = feelPool.CommunityPool.Add(remainder)
		k.SetFeePool(ctx, feelPool)

		if _, _, err := k.bankKeeper.AddCoins(ctx, withdrawAddr, coins); err != nil {
			panic(err)
		}

		outstanding = outstanding.Sub(rewards)

		fmt.Printf("distribute reward[%+v] to delegator[%+v], remainder[%+v] to communityPool\n", rewards, withdrawAddr.String(), remainder)
	}

	if !outstanding.IsZero() {
		feelPool := k.GetFeePool(ctx)
		feelPool.CommunityPool = feelPool.CommunityPool.Add(outstanding)
		k.SetFeePool(ctx, feelPool)
	}
	// set initial historical rewards (period 0) with reference count of 1
	//k.SetValidatorHistoricalRewards(ctx, valAddr, 0, types.NewValidatorHistoricalRewards(sdk.DecCoins{}, 1))

	// set current rewards (starting at period 1)
	k.SetValidatorCurrentRewards(ctx, valAddr, types.NewValidatorCurrentRewards(sdk.DecCoins{}, 1))
	// set accumulated commission
	k.SetValidatorAccumulatedCommission(ctx, valAddr, types.InitialValidatorAccumulatedCommission())
	// set outstanding rewards
	k.SetValidatorOutstandingRewards(ctx, valAddr, sdk.DecCoins{})
}
