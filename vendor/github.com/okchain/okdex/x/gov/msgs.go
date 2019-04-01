package gov

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Governance message types and routes
const (
	TypeMsgDeposit               = "deposit"
	TypeMsgVote                  = "vote"
	TypeMsgSubmitProposal        = "submit_proposal"
	TypeMsgDexListSubmitProposal = "submit_dex_list_proposal"
	TypeMsgDexList               = "dex-list"

	MaxDescriptionLength int = 5000
	MaxTitleLength       int = 140
)

var _, _, _ sdk.Msg = MsgSubmitProposal{}, MsgDeposit{}, MsgVote{}

// MsgSubmitProposal
type MsgSubmitProposal struct {
	Title          string         `json:"title"`           //  Title of the proposal
	Description    string         `json:"description"`     //  Description of the proposal
	ProposalType   ProposalKind   `json:"proposal_type"`   //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       sdk.AccAddress `json:"proposer"`        //  Address of the proposer
	InitialDeposit sdk.Coins      `json:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive.
	Params         Params         `json:"params"`
}

func NewMsgSubmitProposal(title, description string, proposalType ProposalKind, proposer sdk.AccAddress, initialDeposit sdk.Coins, params Params) MsgSubmitProposal {
	return MsgSubmitProposal{
		Title:          title,
		Description:    description,
		ProposalType:   proposalType,
		Proposer:       proposer,
		InitialDeposit: initialDeposit,
		Params:         params,
	}
}

//nolint
func (msg MsgSubmitProposal) Route() string { return RouterKey }
func (msg MsgSubmitProposal) Type() string  { return TypeMsgSubmitProposal }

// Implements Msg.
func (msg MsgSubmitProposal) ValidateBasic() sdk.Error {
	if len(msg.Title) == 0 {
		return ErrInvalidTitle(DefaultCodespace, "No title present in proposal")
	}
	if len(msg.Title) > MaxTitleLength {
		return ErrInvalidTitle(DefaultCodespace, fmt.Sprintf("Proposal title is longer than max length of %d", MaxTitleLength))
	}
	if len(msg.Description) == 0 {
		return ErrInvalidDescription(DefaultCodespace, "No description present in proposal")
	}
	if len(msg.Description) > MaxDescriptionLength {
		return ErrInvalidDescription(DefaultCodespace, fmt.Sprintf("Proposal description is longer than max length of %d", MaxDescriptionLength))
	}
	if !validProposalType(msg.ProposalType) {
		return ErrInvalidProposalType(DefaultCodespace, msg.ProposalType)
	}
	if msg.Proposer.Empty() {
		return sdk.ErrInvalidAddress(msg.Proposer.String())
	}
	if !msg.InitialDeposit.IsValid() {
		return sdk.ErrInvalidCoins(msg.InitialDeposit.String())
	}
	if msg.InitialDeposit.IsAnyNegative() {
		return sdk.ErrInvalidCoins(msg.InitialDeposit.String())
	}
	if msg.ProposalType == ProposalTypeParameterChange {
		if len(msg.Params) == 0 {
			return ErrEmptyParam(DefaultCodespace)
		}
	}
	return nil
}

func (msg MsgSubmitProposal) String() string {
	return fmt.Sprintf("MsgSubmitProposal{%s, %s, %s, %v}", msg.Title, msg.Description, msg.ProposalType, msg.InitialDeposit)
}

// Implements Msg.
func (msg MsgSubmitProposal) GetSignBytes() []byte {
	bz := msgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgSubmitProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proposer}
}

// MsgDexListSubmitProposal
type MsgDexListSubmitProposal struct {
	Title          string         `json:"title"`           //  Title of the proposal
	Description    string         `json:"description"`     //  Description of the proposal
	ProposalType   ProposalKind   `json:"proposal_type"`   //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       sdk.AccAddress `json:"proposer"`        //  Address of the proposer
	InitialDeposit sdk.Coins      `json:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive.
	ListAsset      string         `json:"list_asset"`      //  Symbol of asset listed on Dex.
	QuoteAsset     string         `json:"quote_asset"`     //  Symbol of asset quoted by asset listed on Dex.
	ExpireTime     time.Duration  `json:"expire_time"`     //  Expire time from when proposal is passed.
	InitPrice      sdk.Dec        `json:"init_price"`      //  Init price of asset listed on Dex.
	MaxPriceDigit  uint64         `json:"max_price_digit"` //  Decimal of price
	MaxSizeDigit   uint64         `json:"max_size_digit"`  //  Decimal of trade quantity
	MergeTypes     string         `json:"merge_types"`     //  Level of merge depth
	MinTradeSize   string         `json:"min_trade_size"`
}

func NewMsgDexListSubmitProposal(title, description string, proposalType ProposalKind, proposer sdk.AccAddress, initialDeposit sdk.Coins,
	listAsset, quoteAsset string, expireTime time.Duration, initPrice sdk.Dec, maxPriceDigit, maxSizeDigit uint64, mergeTypes, minTradeSize string) MsgDexListSubmitProposal {
	return MsgDexListSubmitProposal{
		Title:          title,
		Description:    description,
		ProposalType:   proposalType,
		Proposer:       proposer,
		InitialDeposit: initialDeposit,
		ListAsset:      listAsset,
		QuoteAsset:     quoteAsset,
		ExpireTime:     expireTime,
		InitPrice:      initPrice,
		MaxPriceDigit:  maxPriceDigit,
		MaxSizeDigit:   maxSizeDigit,
		MergeTypes:     mergeTypes,
		MinTradeSize:   minTradeSize,
	}
}

//nolint
func (msg MsgDexListSubmitProposal) Route() string { return RouterKey }
func (msg MsgDexListSubmitProposal) Type() string  { return TypeMsgDexListSubmitProposal }

// Implements Msg.
func (msg MsgDexListSubmitProposal) ValidateBasic() sdk.Error {
	if len(msg.Title) == 0 {
		return ErrInvalidTitle(DefaultCodespace, "No title present in proposal")
	}
	if len(msg.Title) > MaxTitleLength {
		return ErrInvalidTitle(DefaultCodespace, fmt.Sprintf("Proposal title is longer than max length of %d", MaxTitleLength))
	}
	if len(msg.Description) == 0 {
		return ErrInvalidDescription(DefaultCodespace, "No description present in proposal")
	}
	if len(msg.Description) > MaxDescriptionLength {
		return ErrInvalidDescription(DefaultCodespace, fmt.Sprintf("Proposal description is longer than max length of %d", MaxDescriptionLength))
	}
	if !validProposalType(msg.ProposalType) {
		return ErrInvalidProposalType(DefaultCodespace, msg.ProposalType)
	}
	if msg.Proposer.Empty() {
		return sdk.ErrInvalidAddress(msg.Proposer.String())
	}
	//if !msg.InitialDeposit.IsValid() {
	//    return sdk.ErrInvalidCoins(msg.InitialDeposit.String())
	//}
	if len(msg.InitialDeposit) != 1 || msg.InitialDeposit[0].Denom != sdk.DefaultBondDenom {
		return sdk.ErrInvalidCoins(fmt.Sprintf("DexList must deposit %s but got %s", sdk.DefaultBondDenom, msg.InitialDeposit.String()))
	}
	if msg.InitialDeposit.IsAnyNegative() {
		return sdk.ErrInvalidCoins(msg.InitialDeposit.String())
	}
	if msg.ListAsset == msg.QuoteAsset {
		return sdk.ErrInvalidCoins(fmt.Sprintf("ListAsset can not equal to QuoteAsset"))
	}
	if msg.QuoteAsset != sdk.DefaultBondDenom {
		return sdk.ErrInvalidCoins(fmt.Sprintf("DexList must quote %s but got %s", sdk.DefaultBondDenom, msg.QuoteAsset))
	}
	if msg.ExpireTime*time.Second < MinExpireTime || msg.ExpireTime*time.Second > MaxExpireTime {
		return ErrInvalidTime(DefaultCodespace, fmt.Sprintf("Dex list expire time must range from %v to %v but got %v", MinExpireTime, MaxExpireTime, msg.ExpireTime*time.Second))
	}
	// TODO:check list asset issued
	return nil
}

func (msg MsgDexListSubmitProposal) String() string {
	return fmt.Sprintf("MsgSubmitProposal{%s, %s, %s, %v, %s, %s, %v, %v}",
		msg.Title, msg.Description, msg.ProposalType, msg.InitialDeposit, msg.ListAsset, msg.QuoteAsset, msg.ExpireTime, msg.InitPrice)
}

// Implements Msg.
func (msg MsgDexListSubmitProposal) GetSignBytes() []byte {
	bz := msgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgDexListSubmitProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proposer}
}

// MsgDeposit
type MsgDeposit struct {
	ProposalID uint64         `json:"proposal_id"` // ID of the proposal
	Depositor  sdk.AccAddress `json:"depositor"`   // Address of the depositor
	Amount     sdk.Coins      `json:"amount"`      // Coins to add to the proposal's deposit
}

func NewMsgDeposit(depositor sdk.AccAddress, proposalID uint64, amount sdk.Coins) MsgDeposit {
	return MsgDeposit{
		ProposalID: proposalID,
		Depositor:  depositor,
		Amount:     amount,
	}
}

// Implements Msg.
// nolint
func (msg MsgDeposit) Route() string { return RouterKey }
func (msg MsgDeposit) Type() string  { return TypeMsgDeposit }

// Implements Msg.
func (msg MsgDeposit) ValidateBasic() sdk.Error {
	if msg.Depositor.Empty() {
		return sdk.ErrInvalidAddress(msg.Depositor.String())
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	if msg.Amount.IsAnyNegative() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	if msg.ProposalID < 0 {
		return ErrUnknownProposal(DefaultCodespace, msg.ProposalID)
	}
	return nil
}

func (msg MsgDeposit) String() string {
	return fmt.Sprintf("MsgDeposit{%s=>%v: %v}", msg.Depositor, msg.ProposalID, msg.Amount)
}

// Implements Msg.
func (msg MsgDeposit) GetSignBytes() []byte {
	bz := msgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Depositor}
}

// MsgVote
type MsgVote struct {
	ProposalID uint64         `json:"proposal_id"` // ID of the proposal
	Voter      sdk.AccAddress `json:"voter"`       //  address of the voter
	Option     VoteOption     `json:"option"`      //  option from OptionSet chosen by the voter
}

func NewMsgVote(voter sdk.AccAddress, proposalID uint64, option VoteOption) MsgVote {
	return MsgVote{
		ProposalID: proposalID,
		Voter:      voter,
		Option:     option,
	}
}

// Implements Msg.
// nolint
func (msg MsgVote) Route() string { return RouterKey }
func (msg MsgVote) Type() string  { return TypeMsgVote }

// Implements Msg.
func (msg MsgVote) ValidateBasic() sdk.Error {
	if msg.Voter.Empty() {
		return sdk.ErrInvalidAddress(msg.Voter.String())
	}
	if msg.ProposalID < 0 {
		return ErrUnknownProposal(DefaultCodespace, msg.ProposalID)
	}
	if !validVoteOption(msg.Option) {
		return ErrInvalidVote(DefaultCodespace, msg.Option)
	}
	return nil
}

func (msg MsgVote) String() string {
	return fmt.Sprintf("MsgVote{%v - %s}", msg.ProposalID, msg.Option)
}

// Implements Msg.
func (msg MsgVote) GetSignBytes() []byte {
	bz := msgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}

type MsgDexList struct {
	Owner      sdk.AccAddress `json:"owner"`
	ProposalID uint64         `json:"proposal-id"`
}

func NewMsgDexList(owner sdk.AccAddress, proposalID uint64) MsgDexList {
	return MsgDexList{
		Owner:      owner,
		ProposalID: proposalID,
	}
}

func (msg MsgDexList) Route() string { return RouterKey }

func (msg MsgDexList) Type() string { return TypeMsgDexList }

func (msg MsgDexList) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if msg.ProposalID < 1 {
		return ErrUnknownProposal(DefaultParamspace, msg.ProposalID)
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDexList) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgDexList) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
