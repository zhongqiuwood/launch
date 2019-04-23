package gov

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ok-chain/okchain/x/gov/tags"
)

// Called every block, process inflation, update validator set
func EndBlocker(ctx sdk.Context, keeper Keeper) sdk.Tags {
	logger := ctx.Logger().With("module", "x/gov")
	resTags := sdk.NewTags()

	// fetch waiting proposals whose voting periods have ended (are passed the block time)
	waitingIterator := keeper.WaitingProposalQueueIterator(ctx, uint64(ctx.BlockHeight()))
	defer waitingIterator.Close()
	for ; waitingIterator.Valid(); waitingIterator.Next() {
		var proposalID uint64

		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(waitingIterator.Value(), &proposalID)
		waitingProposal := keeper.GetProposal(ctx, proposalID)

		err := Execute(ctx, keeper, waitingProposal)
		if err != nil {
			logger.Info(
				fmt.Sprintf("proposal %d (%s) excute failed",
					waitingProposal.GetProposalID(),
					waitingProposal.GetTitle(),
				),
			)
		} else {
			logger.Info(
				fmt.Sprintf("proposal %d (%s) excute successfully",
					waitingProposal.GetProposalID(),
					waitingProposal.GetTitle(),
				),
			)
		}
	}

	inactiveIterator := keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	defer inactiveIterator.Close()
	for ; inactiveIterator.Valid(); inactiveIterator.Next() {
		var proposalID uint64

		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(inactiveIterator.Value(), &proposalID)
		inactiveProposal := keeper.GetProposal(ctx, proposalID)

		keeper.DeleteProposal(ctx, proposalID)
		keeper.DeleteDeposits(ctx, proposalID) // delete any associated deposits (burned)
		if inactiveProposal.GetProposalType() == ProposalTypeDexList && inactiveProposal.(*DexListProposal).BlockHeight != 0 {
			keeper.RefundDexListDeposits(ctx, inactiveProposal.(*DexListProposal).Proposer)
		}
		resTags = resTags.AppendTag(tags.ProposalID, fmt.Sprintf("%d", proposalID))
		resTags = resTags.AppendTag(tags.ProposalResult, tags.ActionProposalDropped)

		keeper.RemoveFromInactiveProposalQueue(ctx, inactiveProposal.GetDepositEndTime(), inactiveProposal.GetProposalID())
		logger.Info(
			fmt.Sprintf("proposal %d (%s) didn't meet minimum deposit of %s (had only %s); deleted",
				inactiveProposal.GetProposalID(),
				inactiveProposal.GetTitle(),
				keeper.GetDepositParams(ctx).MinDeposit,
				inactiveProposal.GetTotalDeposit(),
			),
		)
	}

	// fetch active proposals whose voting periods have ended (are passed the block time)
	activeIterator := keeper.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	defer activeIterator.Close()
	for ; activeIterator.Valid(); activeIterator.Next() {
		var proposalID uint64

		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(activeIterator.Value(), &proposalID)
		activeProposal := keeper.GetProposal(ctx, proposalID)
		passes, tallyResults := tally(ctx, keeper, activeProposal)

		var tagValue string
		if passes {
			if activeProposal.GetProposalType() == ProposalTypeDexList {
				dexListPro := activeProposal.(*DexListProposal)
				dexListPro.SetDexListStartTime(ctx.BlockHeader().Time)
				dexListPro.SetDexListEndTime(dexListPro.GetDexListStartTime().Add(keeper.GetDexListParams(ctx).ExpireTime))
			}
			keeper.RefundDeposits(ctx, activeProposal.GetProposalID())
			activeProposal.SetStatus(StatusPassed)
			tagValue = tags.ActionProposalPassed
			Execute(ctx, keeper, activeProposal)
		} else {
			if activeProposal.GetProposalType() == ProposalTypeDexList && activeProposal.(*DexListProposal).BlockHeight != 0 {
				keeper.RefundDexListDeposits(ctx, activeProposal.(*DexListProposal).Proposer)
			}
			keeper.DeleteDeposits(ctx, activeProposal.GetProposalID())
			activeProposal.SetStatus(StatusRejected)
			tagValue = tags.ActionProposalRejected
		}

		activeProposal.SetFinalTallyResult(tallyResults)
		keeper.SetProposal(ctx, activeProposal)
		keeper.RemoveFromActiveProposalQueue(ctx, activeProposal.GetVotingEndTime(), activeProposal.GetProposalID())

		logger.Info(
			fmt.Sprintf(
				"proposal %d (%s) tallied; passed: %v",
				activeProposal.GetProposalID(), activeProposal.GetTitle(), passes,
			),
		)

		resTags = resTags.AppendTag(tags.ProposalID, fmt.Sprintf("%d", proposalID))
		resTags = resTags.AppendTag(tags.ProposalResult, tagValue)
	}

	return resTags
}
