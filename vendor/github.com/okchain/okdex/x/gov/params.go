package gov

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// DEPOSIT            = 4000
	LOWER_BOUND_AMOUNT = 10
	UPPER_BOUND_AMOUNT = 10000
	// STABLE_NUM         = 1
)

var _ params.ParamSet = (*GovParams)(nil)

//Parameter store key
var (
	KeyMaxDepositPeriod = []byte("MaxDepositPeriod")
	KeyMinDeposit       = []byte("MinDeposit")
	KeyVotingPeriod     = []byte("VotingPeriod")
	KeyQuorum           = []byte("Quorum")
	KeyThreshold        = []byte("Threshold")
	KeyVeto             = []byte("Veto")
)

// ParamTable for gov module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&GovParams{})
}

// mint parameters
type GovParams struct {
	MaxDepositPeriod time.Duration `json:"max_deposit_period"` //  Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months
	MinDeposit       sdk.Coins     `json:"min_deposit"`        //  Minimum deposit for a critical proposal to enter voting period.
	VotingPeriod     time.Duration `json:"voting_period"`      //  Length of the critical voting period.
	Quorum           sdk.Dec       `json:"quorum"`             //
	Threshold        sdk.Dec       `json:"threshold"`          //  Minimum propotion of Yes votes for proposal to pass. Initial value: 0.5
	Veto             sdk.Dec       `json:"veto"`               //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
}

func (p *GovParams) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyMaxDepositPeriod, &p.MaxDepositPeriod},
		{KeyMinDeposit, &p.MinDeposit},
		{KeyVotingPeriod, &p.VotingPeriod},
		{KeyQuorum, &p.Quorum},
		{KeyThreshold, &p.Threshold},
		{KeyVeto, &p.Veto},
	}
}

func (p *GovParams) ValidateKV(key string, value string) (interface{}, sdk.Error) {
	switch key {
	case string(KeyMaxDepositPeriod):
		maxDepositPeriod, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateComplaintRetrospect(complaintRetrospect); err != nil {
		// 	return nil, err
		// }
		return maxDepositPeriod, nil
	case string(KeyMinDeposit):
		minDeposit, err := sdk.ParseCoins(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateArbitrationTimeLimit(arbitrationTimeLimit); err != nil {
		// 	return nil, err
		// }
		return minDeposit, nil
	case string(KeyVotingPeriod):
		votingPeriod, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateComplaintRetrospect(complaintRetrospect); err != nil {
		// 	return nil, err
		// }
		return votingPeriod, nil
	case string(KeyQuorum):
		quorum, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return quorum, nil
	case string(KeyThreshold):
		threshold, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return threshold, nil
	case string(KeyVeto):
		vote, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return vote, nil
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is not found", key))
	}
}

func (p *GovParams) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	case string(KeyMaxDepositPeriod):
		err := cdc.UnmarshalJSON(bytes, &p.MaxDepositPeriod)
		return p.MaxDepositPeriod.String(), err
	case string(KeyMinDeposit):
		err := cdc.UnmarshalJSON(bytes, &p.MinDeposit)
		return p.MinDeposit.String(), err
	case string(KeyVotingPeriod):
		err := cdc.UnmarshalJSON(bytes, &p.VotingPeriod)
		return p.VotingPeriod.String(), err
	case string(KeyQuorum):
		err := cdc.UnmarshalJSON(bytes, &p.Quorum)
		return p.Quorum.String(), err
	case string(KeyThreshold):
		err := cdc.UnmarshalJSON(bytes, &p.Threshold)
		return p.Threshold.String(), err
	case string(KeyVeto):
		err := cdc.UnmarshalJSON(bytes, &p.Veto)
		return p.Veto.String(), err
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

func (p GovParams) String() string {
	return fmt.Sprintf(`Deposit Params:
	Min Deposit:        %s
	Max Deposit Period: %s`, p.MinDeposit, p.MaxDepositPeriod) + "\n" +
		fmt.Sprintf(`Tally Params:
	Quorum:             %s
	Threshold:          %s
	Veto:               %s`,
			p.Quorum, p.Threshold, p.Veto) + "\n" +
		fmt.Sprintf(`Voting Params:
		  Voting Period:      %s`, p.VotingPeriod)
}

func NewGovParams(vp VotingParams, tp TallyParams, dp DepositParams) GovParams {
	return GovParams{
		MaxDepositPeriod: dp.MaxDepositPeriod,
		MinDeposit:       dp.MinDeposit,
		VotingPeriod:     vp.VotingPeriod,
		Quorum:           tp.Quorum,
		Threshold:        tp.Threshold,
		Veto:             tp.Veto,
	}
}

// default minting module parameters
func DefaultParams() GovParams {
	minDepositTokens := sdk.TokensFromTendermintPower(10)
	var minDeposit = sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, minDepositTokens)}

	return GovParams{
		MaxDepositPeriod: DefaultPeriod,
		MinDeposit:       minDeposit,
		VotingPeriod:     DefaultPeriod,
		Quorum:           sdk.NewDecWithPrec(334, 3),
		Threshold:        sdk.NewDecWithPrec(5, 1),
		Veto:             sdk.NewDecWithPrec(334, 3),
	}
}

func validateParams(p GovParams) sdk.Error {
	if err := validateDepositParams(DepositParams{
		MaxDepositPeriod: p.MaxDepositPeriod,
		MinDeposit:       p.MinDeposit,
	}); err != nil {
		return err
	}

	if err := validatorVotingParams(VotingParams{
		VotingPeriod: p.VotingPeriod,
	}); err != nil {
		return err
	}

	if err := validateTallyingParams(TallyParams{
		Quorum:    p.Quorum,
		Threshold: p.Threshold,
		Veto:      p.Veto,
	}); err != nil {
		return err
	}

	return nil
}

//______________________________________________________________________

// get inflation params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) (params GovParams) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// set inflation params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params GovParams) {
	k.paramSpace.SetParamSet(ctx, &params)
}

type DepositParams struct {
	MinDeposit       sdk.Coins
	MaxDepositPeriod time.Duration
}

func (dp DepositParams) String() string {
	return fmt.Sprintf(`Deposit Params:
  Min Deposit:        %s
  Max Deposit Period: %s`, dp.MinDeposit, dp.MaxDepositPeriod)
}

type VotingParams struct {
	VotingPeriod time.Duration `json:"voting_period"` //  Length of the voting period.
}

func (vp VotingParams) String() string {
	return fmt.Sprintf(`Voting Params:
  Voting Period:      %s`, vp.VotingPeriod)
}

// Param around Tallying votes in governance
type TallyParams struct {
	Quorum    sdk.Dec `json:"quorum"`    //  Minimum percentage of total stake needed to vote for a result to be considered valid
	Threshold sdk.Dec `json:"threshold"` //  Minimum propotion of Yes votes for proposal to pass. Initial value: 0.5
	Veto      sdk.Dec `json:"veto"`      //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
}

func (tp TallyParams) String() string {
	return fmt.Sprintf(`Tally Params:
  Quorum:             %s
  Threshold:          %s
  Veto:               %s`,
		tp.Quorum, tp.Threshold, tp.Veto)
}
func validateDepositParams(dp DepositParams) sdk.Error {
	if !dp.MinDeposit.IsValid() {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositDenom, fmt.Sprintf("Governance deposit amount must be a valid sdk.Coins amount, is %s",
			dp.MinDeposit.String()))
	}

	// if dp.MinDeposit[0].Denom != sdk.DefaultBondDenom {
	// 	return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositDenom, fmt.Sprintf("MinDeposit should be %s!", sdk.DefaultBondDenom))
	// }

	// LowerBound, _ := sdk.ParseCoin(fmt.Sprintf("%d%s", LOWER_BOUND_AMOUNT, sdk.DefaultBondDenom))
	// UpperBound, _ := sdk.ParseCoin(fmt.Sprintf("%d%s", UPPER_BOUND_AMOUNT, sdk.DefaultBondDenom))

	// if dp.MinDeposit[0].Amount.LT(LowerBound.Amount) || dp.MinDeposit[0].Amount.GT(UpperBound.Amount) {
	// 	return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositAmount, fmt.Sprintf("MinDepositAmount"+dp.MinDeposit[0].String()+" should be larger than 10iris and less than 10000iris"))
	// }

	// if dp.MaxDepositPeriod < 20*time.Second || dp.MaxDepositPeriod > 3*24*time.Hour {
	// 	return sdk.NewError(params.DefaultCodespace, params.CodeInvalidDepositPeriod, fmt.Sprintf("MaxDepositPeriod (%s) should be between 20s and %d", dp.MaxDepositPeriod.String(), 3*24*time.Hour))
	// }
	return nil
}

func validatorVotingParams(vp VotingParams) sdk.Error {
	if vp.VotingPeriod < 20*time.Second || vp.VotingPeriod > 7*24*time.Hour {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidVotingPeriod, fmt.Sprintf("VotingPeriod (%s) should be between 20s and 1 week", vp.VotingPeriod.String()))
	}
	return nil
}

func validateTallyingParams(tp TallyParams) sdk.Error {
	threshold := tp.Threshold
	if threshold.IsNegative() || threshold.GT(sdk.OneDec()) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxProposalNum, fmt.Sprintf("Governance vote threshold should be positive and less or equal to one, is %s",
			threshold.String()))
	}

	veto := tp.Veto
	if veto.IsNegative() || veto.GT(sdk.OneDec()) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxProposalNum, fmt.Sprintf("Governance vote veto threshold should be positive and less or equal to one, is %s",
			veto.String()))
	}

	// if !tp.Quorum.Equal(sdk.NewDec(STABLE_NUM)) {
	// 	return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxProposalNum, fmt.Sprintf("The num of MaxProposal [%v] can only be %v.", tp.Quorum, sdk.NewDec(STABLE_NUM)))
	// }
	// if tp.Threshold.LTE(sdk.ZeroDec()) || tp.Threshold.GTE(sdk.NewDec(1)) {
	// 	return sdk.NewError(params.DefaultCodespace, params.CodeInvalidThreshold, fmt.Sprintf("Invalid Threshold ( "+tp.Threshold.String()+" ) should be (0,1)"))
	// }
	// if tp.Veto.LTE(sdk.ZeroDec()) || tp.Veto.GTE(sdk.NewDec(1)) {
	// 	return sdk.NewError(params.DefaultCodespace, params.CodeInvalidVeto, fmt.Sprintf("Invalid Veto ( "+tp.Veto.String()+" ) should be (0,1)"))
	// }
	return nil
}
