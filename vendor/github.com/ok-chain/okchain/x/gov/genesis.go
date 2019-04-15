package gov

import (
	"bytes"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Default period for deposits & voting
	DefaultPeriod        time.Duration = 60 * time.Second
	MinExpireTime        time.Duration = 2 * 24 * 3600 * time.Second   // 2 days
	MaxExpireTime        time.Duration = 365 * 24 * 3600 * time.Second // 1 year
	//MaxPriceDigit        uint64        = 8
	//MaxSizeDigit         uint64        = 8
	MaxSumPriceSizeDigit uint64        = 8
	MaxMergeTypeSize     int           = 4
	MaxMergeUint         string        = "100"
	MinMergeFloat        string        = "0.000001"
	MaxMinTradeSizeUint  string        = "100"
	MinMinTradeSizeFloat string        = "0.000001"
	GovTxFee             string        = "0.0125"
)

// GenesisState - all staking state that must be provided at genesis
type GenesisState struct {
	StartingProposalID uint64 `json:"starting_proposal_id"`
	// Deposits           []DepositWithMetadata `json:"deposits"`
	// Votes              []VoteWithMetadata    `json:"votes"`
	Proposals []Proposal `json:"proposals"`
	Params    GovParams  `json:"params"` // params
	// DepositParams      DepositParams         `json:"deposit_params"`
	// VotingParams       VotingParams          `json:"voting_params"`
	// TallyParams        TallyParams           `json:"tally_params"`
}

// DepositWithMetadata (just for genesis)
type DepositWithMetadata struct {
	ProposalID uint64  `json:"proposal_id"`
	Deposit    Deposit `json:"deposit"`
}

// VoteWithMetadata (just for genesis)
type VoteWithMetadata struct {
	ProposalID uint64 `json:"proposal_id"`
	Vote       Vote   `json:"vote"`
}

func NewGenesisState(startingProposalID uint64, params GovParams) GenesisState {
	return GenesisState{
		StartingProposalID: startingProposalID,
		// DepositParams:      dp,
		// VotingParams:       vp,
		// TallyParams:        tp,
		Params: params,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		StartingProposalID: 1,
		Params:             DefaultParams(),
		// DepositParams: DepositParams{
		// 	MinDeposit:       sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, minDepositTokens)},
		// 	MaxDepositPeriod: DefaultPeriod,
		// },
		// VotingParams: VotingParams{
		// 	VotingPeriod: DefaultPeriod,
		// },
		// TallyParams: TallyParams{
		// 	Quorum:    sdk.NewDecWithPrec(334, 3),
		// 	Threshold: sdk.NewDecWithPrec(5, 1),
		// 	Veto:      sdk.NewDecWithPrec(334, 3),
		// },
	}
}

// Checks whether 2 GenesisState structs are equivalent.
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := msgCdc.MustMarshalBinaryBare(data)
	b2 := msgCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// Returns if a GenesisState is empty or has data in it
func (data GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return data.Equal(emptyGenState)
}

// ValidateGenesis
func ValidateGenesis(data GenesisState) error {
	err := validateParams(data.Params)
	if err != nil {
		return err
	}
	if data.Params.MaxDepositPeriod > data.Params.VotingPeriod {
		return fmt.Errorf("Governance deposit period should be less than or equal to the voting period (%ds), is %ds",
			data.Params.VotingPeriod, data.Params.MaxDepositPeriod)
	}
	return nil
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	err := k.setInitialProposalID(ctx, data.StartingProposalID)
	if err != nil {
		// TODO: Handle this with #870
		panic(err)
	}
	// k.setDepositParams(ctx, data.DepositParams)
	// k.setVotingParams(ctx, data.VotingParams)
	// k.setTallyParams(ctx, data.TallyParams)
	// for _, deposit := range data.Deposits {
	// 	k.setDeposit(ctx, deposit.ProposalID, deposit.Deposit.Depositor, deposit.Deposit)
	// }
	// for _, vote := range data.Votes {
	// 	k.setVote(ctx, vote.ProposalID, vote.Vote.Voter, vote.Vote)
	// }
	for _, proposal := range data.Proposals {
		switch proposal.GetStatus() {
		case StatusDepositPeriod:
			k.InsertInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposal.GetProposalID())
		case StatusVotingPeriod:
			k.InsertActiveProposalQueue(ctx, proposal.GetVotingEndTime(), proposal.GetProposalID())
		}
		k.SetProposal(ctx, proposal)
	}
	k.SetParams(ctx, data.Params)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	startingProposalID, _ := k.peekCurrentProposalID(ctx)
	depositParams := k.GetDepositParams(ctx)
	dexListDepositParams := k.GetDexListDepositParams(ctx)
	votingParams := k.GetVotingParams(ctx)
	dexListVotingParams := k.GetDexListVotingParams(ctx)
	dexListParams := k.GetDexListParams(ctx)
	tallyParams := k.GetTallyParams(ctx)
	var deposits []DepositWithMetadata
	var votes []VoteWithMetadata
	proposals := k.GetProposalsFiltered(ctx, nil, nil, StatusNil, 0)
	for _, proposal := range proposals {
		proposalID := proposal.GetProposalID()
		depositsIterator := k.GetDeposits(ctx, proposalID)
		defer depositsIterator.Close()
		for ; depositsIterator.Valid(); depositsIterator.Next() {
			var deposit Deposit
			k.cdc.MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), &deposit)
			deposits = append(deposits, DepositWithMetadata{proposalID, deposit})
		}
		votesIterator := k.GetVotes(ctx, proposalID)
		defer votesIterator.Close()
		for ; votesIterator.Valid(); votesIterator.Next() {
			var vote Vote
			k.cdc.MustUnmarshalBinaryLengthPrefixed(votesIterator.Value(), &vote)
			votes = append(votes, VoteWithMetadata{proposalID, vote})
		}
	}

	return GenesisState{
		StartingProposalID: startingProposalID,
		// Deposits:           deposits,
		// Votes:              votes,
		Proposals: proposals,
		Params: GovParams{
			depositParams.MaxDepositPeriod,
			depositParams.MinDeposit,
			votingParams.VotingPeriod,
			dexListDepositParams.MaxDepositPeriod,
			dexListDepositParams.MinDeposit,
			dexListVotingParams.VotingPeriod,
			dexListVotingParams.VotingFee,
			dexListParams.MaxBlockHeight,
			dexListParams.Fee,
			dexListParams.ExpireTime,
			tallyParams.Quorum,
			tallyParams.Threshold,
			tallyParams.Veto,
			0,
		},
	}
}
