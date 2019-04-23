package gov

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Insert string = "insert"
	Update string = "update"
)

// Proposal interface
type Proposal interface {
	GetProposalID() uint64
	SetProposalID(uint64)

	GetTitle() string
	SetTitle(string)

	GetDescription() string
	SetDescription(string)

	GetProposalType() ProposalKind
	SetProposalType(ProposalKind)

	GetStatus() ProposalStatus
	SetStatus(ProposalStatus)

	GetFinalTallyResult() TallyResult
	SetFinalTallyResult(TallyResult)

	GetSubmitTime() time.Time
	SetSubmitTime(time.Time)

	GetDepositEndTime() time.Time
	SetDepositEndTime(time.Time)

	GetTotalDeposit() sdk.DecCoins
	SetTotalDeposit(sdk.DecCoins)

	GetVotingStartTime() time.Time
	SetVotingStartTime(time.Time)

	GetVotingEndTime() time.Time
	SetVotingEndTime(time.Time)

	String() string
}

// Proposals is an array of proposal
type Proposals []Proposal

func (p Proposals) String() string {
	out := "ID - (Status) [Type] Title\n"
	for _, prop := range p {
		out += fmt.Sprintf("%d - (%s) [%s] %s\n",
			prop.GetProposalID(), prop.GetStatus(),
			prop.GetProposalType(), prop.GetTitle())
	}
	return strings.TrimSpace(out)
}

// checks if two proposals are equal
func ProposalEqual(proposalA Proposal, proposalB Proposal) bool {
	if proposalA.GetProposalID() == proposalB.GetProposalID() &&
		proposalA.GetTitle() == proposalB.GetTitle() &&
		proposalA.GetDescription() == proposalB.GetDescription() &&
		proposalA.GetProposalType() == proposalB.GetProposalType() &&
		proposalA.GetStatus() == proposalB.GetStatus() &&
		proposalA.GetFinalTallyResult().Equals(proposalB.GetFinalTallyResult()) &&
		proposalA.GetSubmitTime().Equal(proposalB.GetSubmitTime()) &&
		proposalA.GetDepositEndTime().Equal(proposalB.GetDepositEndTime()) &&
		proposalA.GetTotalDeposit().IsEqual(proposalB.GetTotalDeposit()) &&
		proposalA.GetVotingStartTime().Equal(proposalB.GetVotingStartTime()) &&
		proposalA.GetVotingEndTime().Equal(proposalB.GetVotingEndTime()) {
		return true
	}
	return false
}

// Basic Proposals
type BasicProposal struct {
	ProposalID   uint64       `json:"proposal_id"`   //  ID of the proposal
	Title        string       `json:"title"`         //  Title of the proposal
	Description  string       `json:"description"`   //  Description of the proposal
	ProposalType ProposalKind `json:"proposal_type"` //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}

	Status           ProposalStatus `json:"proposal_status"` //  Status of the Proposal {Pending, Active, Passed, Rejected}
	FinalTallyResult TallyResult    `json:"tally_result"`    //  Result of Tallys

	SubmitTime     time.Time    `json:"submit_time"`      //  Time of the block where TxGovSubmitProposal was included
	DepositEndTime time.Time    `json:"deposit_end_time"` // Time that the Proposal would expire if deposit amount isn't met
	TotalDeposit   sdk.DecCoins `json:"total_deposit"`    //  Current deposit on this proposal. Initial value is set at InitialDeposit

	VotingStartTime time.Time `json:"voting_start_time"` //  Time of the block where MinDeposit was reached. -1 if MinDeposit is not reached
	VotingEndTime   time.Time `json:"voting_end_time"`   // Time that the VotingPeriod for this proposal will end and votes will be tallied
}

// Implements Proposal Interface
var _ Proposal = (*BasicProposal)(nil)

// nolint
func (tp BasicProposal) GetProposalID() uint64                      { return tp.ProposalID }
func (tp *BasicProposal) SetProposalID(proposalID uint64)           { tp.ProposalID = proposalID }
func (tp BasicProposal) GetTitle() string                           { return tp.Title }
func (tp *BasicProposal) SetTitle(title string)                     { tp.Title = title }
func (tp BasicProposal) GetDescription() string                     { return tp.Description }
func (tp *BasicProposal) SetDescription(description string)         { tp.Description = description }
func (tp BasicProposal) GetProposalType() ProposalKind              { return tp.ProposalType }
func (tp *BasicProposal) SetProposalType(proposalType ProposalKind) { tp.ProposalType = proposalType }
func (tp BasicProposal) GetStatus() ProposalStatus                  { return tp.Status }
func (tp *BasicProposal) SetStatus(status ProposalStatus)           { tp.Status = status }
func (tp BasicProposal) GetFinalTallyResult() TallyResult           { return tp.FinalTallyResult }
func (tp *BasicProposal) SetFinalTallyResult(tallyResult TallyResult) {
	tp.FinalTallyResult = tallyResult
}
func (tp BasicProposal) GetSubmitTime() time.Time            { return tp.SubmitTime }
func (tp *BasicProposal) SetSubmitTime(submitTime time.Time) { tp.SubmitTime = submitTime }
func (tp BasicProposal) GetDepositEndTime() time.Time        { return tp.DepositEndTime }
func (tp *BasicProposal) SetDepositEndTime(depositEndTime time.Time) {
	tp.DepositEndTime = depositEndTime
}
func (tp BasicProposal) GetTotalDeposit() sdk.DecCoins              { return tp.TotalDeposit }
func (tp *BasicProposal) SetTotalDeposit(totalDeposit sdk.DecCoins) { tp.TotalDeposit = totalDeposit }
func (tp BasicProposal) GetVotingStartTime() time.Time              { return tp.VotingStartTime }
func (tp *BasicProposal) SetVotingStartTime(votingStartTime time.Time) {
	tp.VotingStartTime = votingStartTime
}
func (tp BasicProposal) GetVotingEndTime() time.Time { return tp.VotingEndTime }
func (tp *BasicProposal) SetVotingEndTime(votingEndTime time.Time) {
	tp.VotingEndTime = votingEndTime
}

func (tp BasicProposal) String() string {
	return fmt.Sprintf(`Proposal %d:
  Title:              %s
  Type:               %s
  Status:             %s
  Submit Time:        %s
  Deposit End Time:   %s
  Total Deposit:      %s
  Voting Start Time:  %s
  Voting End Time:    %s`, tp.ProposalID, tp.Title, tp.ProposalType,
		tp.Status, tp.SubmitTime, tp.DepositEndTime,
		tp.TotalDeposit, tp.VotingStartTime, tp.VotingEndTime)
}

// Text Proposals
type TextProposal struct {
	BasicProposal
}

// Implements Proposal Interface
var _ Proposal = (*TextProposal)(nil)

// DexList Proposals
type DexListProposal struct {
	BasicProposal
	Proposer      sdk.AccAddress `json:"proposer"`    //  Proposer of proposal
	ListAsset     string         `json:"list_asset"`  //  Symbol of asset listed on Dex.
	QuoteAsset    string         `json:"quote_asset"` //  Symbol of asset quoted by asset listed on Dex.
	InitPrice     sdk.Dec        `json:"init_price"`  //  Init price of asset listed on Dex.
	BlockHeight   uint64         `json:"block_height"`
	MaxPriceDigit uint64         `json:"max_price_digit"` //  Decimal of price
	MaxSizeDigit  uint64         `json:"max_size_digit"`  //  Decimal of trade quantity
	MinTradeSize  string         `json:"min_trade_size"`

	DexListStartTime time.Time `json:"dex_list_start_time"`
	DexListEndTime   time.Time `json:"dex_list_end_time"`
}

// Implements Proposal Interface
var _ Proposal = (*DexListProposal)(nil)

func (tp DexListProposal) String() string {
	return fmt.Sprintf(`Proposal %d:
  Title:               %s
  Type:                %s
  Proposer:            %s
  Status:              %s
  Submit Time:         %s
  Deposit End Time:    %s
  Total Deposit:       %s
  Voting Start Time:   %s
  Voting End Time:     %s
  ListAsset            %s
  QuoteAsset           %s
  InitPrice            %s
  BlockHeight          %d
  MaxPriceDigit        %d
  MaxSizeDigit         %d
  MinTradeSize         %s
  Dex List Start Time: %s
  Dex List Start Time: %s`, tp.ProposalID, tp.Title, tp.ProposalType, tp.Proposer,
		tp.Status, tp.SubmitTime, tp.DepositEndTime,
		tp.TotalDeposit, tp.VotingStartTime, tp.VotingEndTime, tp.ListAsset, tp.QuoteAsset, tp.InitPrice,
		tp.BlockHeight, tp.MaxPriceDigit, tp.MaxSizeDigit, tp.MinTradeSize, tp.DexListStartTime, tp.DexListEndTime)
}

func (tp *DexListProposal) SetDexListStartTime(startTime time.Time) {
	tp.DexListStartTime = startTime
}
func (tp *DexListProposal) GetDexListStartTime() time.Time { return tp.DexListStartTime }
func (tp *DexListProposal) SetDexListEndTime(endTime time.Time) {
	tp.DexListEndTime = endTime
}
func (tp *DexListProposal) GetDexListEndTime() time.Time { return tp.DexListEndTime }

type Param struct {
	Subspace string `json:"subspace"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

type Params []Param

// Implements Proposal Interface
var _ Proposal = (*ParameterProposal)(nil)

type ParameterProposal struct {
	BasicProposal
	Params Params `json:"params"`
	Height int64  `json:"height"`
}

// ProposalQueue
type ProposalQueue []uint64

// ProposalKind

// Type that represents Proposal Type as a byte
type ProposalKind byte

//nolint
const (
	ProposalTypeNil             ProposalKind = 0x00
	ProposalTypeText            ProposalKind = 0x01
	ProposalTypeParameterChange ProposalKind = 0x02
	ProposalTypeSoftwareUpgrade ProposalKind = 0x03
	ProposalTypeDexList         ProposalKind = 0x04
)

// String to proposalType byte. Returns 0xff if invalid.
func ProposalTypeFromString(str string) (ProposalKind, error) {
	switch str {
	case "Text":
		return ProposalTypeText, nil
	case "ParameterChange":
		return ProposalTypeParameterChange, nil
	case "SoftwareUpgrade":
		return ProposalTypeSoftwareUpgrade, nil
	case "DexList":
		return ProposalTypeDexList, nil
	default:
		return ProposalKind(0xff), fmt.Errorf("'%s' is not a valid proposal type", str)
	}
}

// is defined ProposalType?
func validProposalType(pt ProposalKind) bool {
	if pt == ProposalTypeText ||
		pt == ProposalTypeParameterChange ||
		pt == ProposalTypeSoftwareUpgrade ||
		pt == ProposalTypeDexList {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (pt ProposalKind) Marshal() ([]byte, error) {
	return []byte{byte(pt)}, nil
}

// Unmarshal needed for protobuf compatibility
func (pt *ProposalKind) Unmarshal(data []byte) error {
	*pt = ProposalKind(data[0])
	return nil
}

// Marshals to JSON using string
func (pt ProposalKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(pt.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (pt *ProposalKind) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := ProposalTypeFromString(s)
	if err != nil {
		return err
	}
	*pt = bz2
	return nil
}

// Turns VoteOption byte to String
func (pt ProposalKind) String() string {
	switch pt {
	case ProposalTypeText:
		return "Text"
	case ProposalTypeParameterChange:
		return "ParameterChange"
	case ProposalTypeSoftwareUpgrade:
		return "SoftwareUpgrade"
	case ProposalTypeDexList:
		return "DexList"
	default:
		return ""
	}
}

// For Printf / Sprintf, returns bech32 when using %s
// nolint: errcheck
func (pt ProposalKind) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(pt.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(pt))))
	}
}

// ProposalStatus

// Type that represents Proposal Status as a byte
type ProposalStatus byte

//nolint
const (
	StatusNil           ProposalStatus = 0x00
	StatusDepositPeriod ProposalStatus = 0x01
	StatusVotingPeriod  ProposalStatus = 0x02
	StatusPassed        ProposalStatus = 0x03
	StatusRejected      ProposalStatus = 0x04
)

// ProposalStatusToString turns a string into a ProposalStatus
func ProposalStatusFromString(str string) (ProposalStatus, error) {
	switch str {
	case "DepositPeriod":
		return StatusDepositPeriod, nil
	case "VotingPeriod":
		return StatusVotingPeriod, nil
	case "Passed":
		return StatusPassed, nil
	case "Rejected":
		return StatusRejected, nil
	case "":
		return StatusNil, nil
	default:
		return ProposalStatus(0xff), fmt.Errorf("'%s' is not a valid proposal status", str)
	}
}

// is defined ProposalType?
func validProposalStatus(status ProposalStatus) bool {
	if status == StatusDepositPeriod ||
		status == StatusVotingPeriod ||
		status == StatusPassed ||
		status == StatusRejected {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (status ProposalStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *ProposalStatus) Unmarshal(data []byte) error {
	*status = ProposalStatus(data[0])
	return nil
}

// Marshals to JSON using string
func (status ProposalStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (status *ProposalStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := ProposalStatusFromString(s)
	if err != nil {
		return err
	}
	*status = bz2
	return nil
}

// Turns VoteStatus byte to String
func (status ProposalStatus) String() string {
	switch status {
	case StatusDepositPeriod:
		return "DepositPeriod"
	case StatusVotingPeriod:
		return "VotingPeriod"
	case StatusPassed:
		return "Passed"
	case StatusRejected:
		return "Rejected"
	default:
		return ""
	}
}

// For Printf / Sprintf, returns bech32 when using %s
// nolint: errcheck
func (status ProposalStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}

// Tally Results
type TallyResult struct {
	Yes        sdk.Int `json:"yes"`
	Abstain    sdk.Int `json:"abstain"`
	No         sdk.Int `json:"no"`
	NoWithVeto sdk.Int `json:"no_with_veto"`
}

func NewTallyResult(yes, abstain, no, noWithVeto sdk.Int) TallyResult {
	return TallyResult{
		Yes:        yes,
		Abstain:    abstain,
		No:         no,
		NoWithVeto: noWithVeto,
	}
}

func NewTallyResultFromMap(results map[VoteOption]sdk.Dec) TallyResult {
	return TallyResult{
		Yes:        results[OptionYes].TruncateInt(),
		Abstain:    results[OptionAbstain].TruncateInt(),
		No:         results[OptionNo].TruncateInt(),
		NoWithVeto: results[OptionNoWithVeto].TruncateInt(),
	}
}

// checks if two proposals are equal
func EmptyTallyResult() TallyResult {
	return TallyResult{
		Yes:        sdk.ZeroInt(),
		Abstain:    sdk.ZeroInt(),
		No:         sdk.ZeroInt(),
		NoWithVeto: sdk.ZeroInt(),
	}
}

// checks if two proposals are equal
func (tr TallyResult) Equals(comp TallyResult) bool {
	return (tr.Yes.Equal(comp.Yes) &&
		tr.Abstain.Equal(comp.Abstain) &&
		tr.No.Equal(comp.No) &&
		tr.NoWithVeto.Equal(comp.NoWithVeto))
}

func (tr TallyResult) String() string {
	return fmt.Sprintf(`Tally Result:
  Yes:        %s
  Abstain:    %s
  No:         %s
  NoWithVeto: %s`, tr.Yes, tr.Abstain, tr.No, tr.NoWithVeto)
}
