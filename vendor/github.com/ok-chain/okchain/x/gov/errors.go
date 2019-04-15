//nolint
package gov

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeUnknownProposal         sdk.CodeType = 1
	CodeInactiveProposal        sdk.CodeType = 2
	CodeAlreadyActiveProposal   sdk.CodeType = 3
	CodeAlreadyFinishedProposal sdk.CodeType = 4
	CodeAddressNotStaked        sdk.CodeType = 5
	CodeInvalidTitle            sdk.CodeType = 6
	CodeInvalidDescription      sdk.CodeType = 7
	CodeInvalidProposalType     sdk.CodeType = 8
	CodeInvalidVote             sdk.CodeType = 9
	CodeInvalidGenesis          sdk.CodeType = 10
	CodeInvalidProposalStatus   sdk.CodeType = 11
	CodeInvalidTime             sdk.CodeType = 12
	CodeInvalidAsset            sdk.CodeType = 13
	CodeEmptyParam              sdk.CodeType = 14
	CodeInvalidParam            sdk.CodeType = 15
	CodeInvalidPriceDigit       sdk.CodeType = 16
	CodeInvalidSizeDigit        sdk.CodeType = 17
	CodeInvalidMergeTypes       sdk.CodeType = 18
	CodeInvalidMinTradeSize     sdk.CodeType = 19
	CodeInvalidHeight           sdk.CodeType = 20
	CodeInvalidMaxHeightPeriod  sdk.CodeType = 21
)

// Error constructors

func ErrUnknownProposal(codespace sdk.CodespaceType, proposalID uint64) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownProposal, fmt.Sprintf("Unknown proposal with id %d", proposalID))
}

func ErrInactiveProposal(codespace sdk.CodespaceType, proposalID uint64) sdk.Error {
	return sdk.NewError(codespace, CodeInactiveProposal, fmt.Sprintf("Inactive proposal with id %d", proposalID))
}

func ErrAlreadyActiveProposal(codespace sdk.CodespaceType, proposalID uint64) sdk.Error {
	return sdk.NewError(codespace, CodeAlreadyActiveProposal, fmt.Sprintf("Proposal %d has been already active", proposalID))
}

func ErrAlreadyFinishedProposal(codespace sdk.CodespaceType, proposalID uint64) sdk.Error {
	return sdk.NewError(codespace, CodeAlreadyFinishedProposal, fmt.Sprintf("Proposal %d has already passed its voting period", proposalID))
}

func ErrAddressNotStaked(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeAddressNotStaked, fmt.Sprintf("Address %s is not staked and is thus ineligible to vote", address))
}

func ErrInvalidTitle(codespace sdk.CodespaceType, errorMsg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTitle, errorMsg)
}

func ErrInvalidDescription(codespace sdk.CodespaceType, errorMsg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDescription, errorMsg)
}

func ErrInvalidProposalType(codespace sdk.CodespaceType, proposalType ProposalKind) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProposalType, fmt.Sprintf("Proposal Type '%s' is not valid", proposalType))
}

func ErrInvalidVote(codespace sdk.CodespaceType, voteOption VoteOption) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidVote, fmt.Sprintf("'%v' is not a valid voting option", voteOption))
}

func ErrInvalidGenesis(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidVote, msg)
}

func ErrInvalidTime(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTime, msg)
}

func ErrEmptyParam(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyParam, fmt.Sprintf("Params can't be empty"))
}

func ErrInvalidParam(codespace sdk.CodespaceType, str string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidParam, fmt.Sprintf("%s Params don't support the ParameterChange.", str))
}

func ErrInvalHeight(codespace sdk.CodespaceType, h, ch, max int64) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidHeight, fmt.Sprintf("Height %d must be greater than current block height %d and less than %d + %d.", h, ch, ch, max))
}

func ErrInvalMaxHeightPeriod(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMaxHeightPeriod, fmt.Sprintf("Height must be greater than 0."))
}
