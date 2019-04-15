package gov

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/ok-chain/okchain/x/token"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ok-chain/okchain/x/gov/tags"
)

// Handle all "gov" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDeposit:
			return handleMsgDeposit(ctx, keeper, msg)
		case MsgSubmitProposal:
			return handleMsgSubmitProposal(ctx, keeper, msg)
		case MsgVote:
			return handleMsgVote(ctx, keeper, msg)
		case MsgDexListSubmitProposal:
			return handleMsgDexListSubmitProposal(ctx, keeper, msg)
		case MsgDexList:
			return handleMsgDexList(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized gov msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSubmitProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitProposal) sdk.Result {
	err := keeper.AddCollectedFees(ctx, keeper.GetGovTxFee(), msg.Proposer)
	if err != nil {
		return err.Result()
	}
	if msg.ProposalType == ProposalTypeParameterChange {
		curHeight := ctx.BlockHeight()
		maxHeight := keeper.GetParams(ctx).MaxBlockHeightPeriod
		if maxHeight == 0 {
			maxHeight = math.MaxInt64 - msg.Height
		}
		if msg.Height <= curHeight || msg.Height > curHeight+maxHeight {
			return ErrInvalHeight(DefaultCodespace, msg.Height, curHeight, maxHeight).Result()
		}
		for _, param := range msg.Params {
			if _, ok := keeper.paramsKeeper.GetParamSet(param.Subspace); ok {
				// if _, err := p.Validate(param.Key, param.Value); err != nil {
				// 	return err.Result()
				// }
			} else {
				return ErrInvalidParam(DefaultCodespace, param.Subspace).Result()
			}
		}
	}
	proposal := keeper.NewProposal(ctx, msg.Title, msg.Description, msg.ProposalType, msg.Params, msg.Height)
	proposalID := proposal.GetProposalID()
	proposalIDStr := fmt.Sprintf("%d", proposalID)

	err, votingStarted := keeper.AddDeposit(ctx, proposalID, msg.Proposer, msg.InitialDeposit)
	if err != nil {
		return err.Result()
	}

	resTags := sdk.NewTags(
		tags.Proposer, []byte(msg.Proposer.String()),
		tags.ProposalID, proposalIDStr,
	)

	if msg.ProposalType == ProposalTypeParameterChange {
		var paramBytes []byte
		paramBytes, _ = json.Marshal(proposal.(*ParameterProposal).Params)
		resTags.AppendTag(tags.Param, string(paramBytes))
	}

	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDStr)
	}

	return sdk.Result{
		Data: keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID),
		Tags: resTags,
	}
}

func handleMsgDexListSubmitProposal(ctx sdk.Context, keeper Keeper, msg MsgDexListSubmitProposal) sdk.Result {
	err := keeper.AddCollectedFees(ctx, keeper.GetGovTxFee(), msg.Proposer)
	if err != nil {
		return err.Result()
	}
	// check asset is issued.
	//if keeper.tokenKeeper == nil {
	//	return sdk.NewError(DefaultParamspace, CodeInvalidGenesis, fmt.Sprintf("tokenKeeper in gov keeper is nil")).Result()
	//}
	if keeper.tokenKeeper.GetTokenInfo(ctx, msg.ListAsset).Symbol != msg.ListAsset {
		return sdk.NewError(DefaultParamspace, CodeInvalidAsset, fmt.Sprintf("asset %s has not been issued", msg.ListAsset)).Result()
	}

	maxBlockHeight := keeper.GetDexListParams(ctx).MaxBlockHeight
	if msg.BlockHeight > maxBlockHeight {
		return sdk.NewError(DefaultParamspace, CodeInvalidParam, fmt.Sprintf("dex list block height can not be greater than %d", maxBlockHeight)).Result()
	}

	proposal := keeper.NewDexListProposal(ctx, msg)
	proposalID := proposal.GetProposalID()
	proposalIDStr := fmt.Sprintf("%d", proposalID)

	if msg.BlockHeight != 0 {
		_, err := keeper.ck.SendCoins(ctx, msg.Proposer, DexListDepositedCoinsAccAddr, NewCoinsFromDecCoins(keeper.GetDexListParams(ctx).Fee))
		if err != nil {
			return err.Result()
		}
	}

	err, votingStarted := keeper.AddDeposit(ctx, proposalID, msg.Proposer, msg.InitialDeposit)
	if err != nil {
		return err.Result()
	}

	resTags := sdk.NewTags(
		tags.Proposer, []byte(msg.Proposer.String()),
		tags.ProposalID, proposalIDStr,
	)

	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDStr)
	}

	return sdk.Result{
		Data: keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID),
		Tags: resTags,
	}
}

func handleMsgDeposit(ctx sdk.Context, keeper Keeper, msg MsgDeposit) sdk.Result {
	err := keeper.AddCollectedFees(ctx, keeper.GetGovTxFee(), msg.Depositor)
	if err != nil {
		return err.Result()
	}
	err, votingStarted := keeper.AddDeposit(ctx, msg.ProposalID, msg.Depositor, msg.Amount)
	if err != nil {
		return err.Result()
	}

	proposalIDStr := fmt.Sprintf("%d", msg.ProposalID)
	resTags := sdk.NewTags(
		tags.Depositor, []byte(msg.Depositor.String()),
		tags.ProposalID, proposalIDStr,
	)

	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDStr)
	}

	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgVote(ctx sdk.Context, keeper Keeper, msg MsgVote) sdk.Result {
	err := keeper.AddCollectedFees(ctx, keeper.GetGovTxFee(), msg.Voter)
	if err != nil {
		// TODO: panic should not happen
		return err.Result()
	}
	err = keeper.AddVote(ctx, msg.ProposalID, msg.Voter, msg.Option)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Voter, msg.Voter.String(),
			tags.ProposalID, fmt.Sprintf("%d", msg.ProposalID),
		),
	}
}

func handleMsgDexList(ctx sdk.Context, keeper Keeper, msg MsgDexList) sdk.Result {
	err := keeper.AddCollectedFees(ctx, keeper.GetGovTxFee(), msg.Owner)
	if err != nil {
		// TODO: panic should not happen
		return err.Result()
	}
	proposal := keeper.GetProposal(ctx, msg.ProposalID)
	dexListProposal, ok := proposal.(*DexListProposal)
	if !ok {
		return sdk.NewError(DefaultParamspace, CodeInvalidProposalType, fmt.Sprintf("proposal is not DexList proposal")).Result()
	}
	if dexListProposal.BlockHeight != 0 {
		return sdk.NewError(DefaultParamspace, CodeInvalidProposalType, fmt.Sprintf("proposal must be excuted at block %d", dexListProposal.BlockHeight)).Result()
	}
	// check owner
	if !bytes.Equal(dexListProposal.Proposer.Bytes(), msg.Owner.Bytes()) {
		return sdk.ErrUnauthorized("Not the proposal's owner").Result()
	}
	if dexListProposal.Status != StatusPassed {
		return ErrInactiveProposal(DefaultParamspace, dexListProposal.ProposalID).Result()
	}
	if ctx.BlockHeader().Time.After(dexListProposal.GetDexListEndTime()) {
		return sdk.NewError(DefaultParamspace, CodeInvalidProposalType, fmt.Sprintf("proposal has expired time %v", dexListProposal.DexListEndTime)).Result()
	}
	return saveTokenPair(ctx, keeper, dexListProposal)
}

func saveTokenPair(ctx sdk.Context, keeper Keeper, dexListProposal *DexListProposal) sdk.Result {
	tokenPair := token.TokenPair{
		BaseAssetSymbol:  dexListProposal.ListAsset,
		QuoteAssetSymbol: dexListProposal.QuoteAsset,
		InitPrice:        dexListProposal.InitPrice,
		MaxPriceDigit:    int64(dexListProposal.MaxPriceDigit),
		MaxQuantityDigit: int64(dexListProposal.MaxSizeDigit),
		//MergeTypes:       dexListProposal.MergeTypes,
		MinQuantity: sdk.MustNewDecFromStr(dexListProposal.MinTradeSize),
	}
	keeper.tokenKeeper.SaveTokenPair(ctx, tokenPair)
	return sdk.Result{
		Tags: sdk.NewTags(
			"list-asset", tokenPair.BaseAssetSymbol,
			"quote-asset", tokenPair.QuoteAssetSymbol,
			"init-price", tokenPair.InitPrice.String(),
			"max-price-digit", strconv.FormatInt(tokenPair.MaxPriceDigit, 10),
			"max-size-digit", strconv.FormatInt(tokenPair.MaxQuantityDigit, 10),
			//"merge-types", tokenPair.MergeTypes,
			"min-trade-size", tokenPair.MinQuantity.String(),
		),
	}
}