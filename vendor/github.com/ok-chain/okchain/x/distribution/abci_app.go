package distribution

import (
	"github.com/ok-chain/okchain/x/common"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ok-chain/okchain/x/distribution/keeper"
	"github.com/ok-chain/okchain/x/distribution/types"
)

// set the proposer for determining distribution during endblock
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper, feeKeeper types.FeeCollectionKeeper) {
	// determine the total power signing the block
	var totalPower, sumPrecommitPower int64
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		totalPower += voteInfo.Validator.Power
		if voteInfo.SignedLastBlock {
			sumPrecommitPower += voteInfo.Validator.Power
		}
	}

	// TODO this is Tendermint-dependent
	// ref https://github.com/cosmos/cosmos-sdk/issues/3095
	if ctx.BlockHeight() > 1 {
		previousProposer := k.GetPreviousProposerConsAddr(ctx)
		k.AllocateTokens(ctx, sumPrecommitPower, totalPower, previousProposer, req.LastCommitInfo.GetVotes())
	}

	// record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	k.SetPreviousProposerConsAddr(ctx, consAddr)

	if _, ok := common.IsEpochEnd(ctx); ok {
		k.DistributeAllRewards(ctx)
		k.SetValidatorsSnapshoot(ctx) //Snapshoot validator set, and it's belonged delegations
	}
}
