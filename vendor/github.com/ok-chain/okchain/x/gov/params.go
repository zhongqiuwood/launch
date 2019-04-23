package gov

import (
	"fmt"
	"strconv"
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
	KeyMaxDepositPeriod        = []byte("MaxDepositPeriod")
	KeyMinDeposit              = []byte("MinDeposit")
	KeyVotingPeriod            = []byte("VotingPeriod")
	KeyDexListMaxDepositPeriod = []byte("DexListMaxDepositPeriod")
	KeyDexListMinDeposit       = []byte("DexListMinDeposit")
	KeyDexListVotingPeriod     = []byte("DexListVotingPeriod")
	KeyDexListVoteFee          = []byte("DexListVoteFee")
	KeyDexListMaxBlockHeight   = []byte("DexListMaxBlockHeight")
	KeyDexListFee              = []byte("DexListFee")
	KeyDexListExpireTime       = []byte("DexListExpireTime")
	KeyQuorum                  = []byte("Quorum")
	KeyThreshold               = []byte("Threshold")
	KeyVeto                    = []byte("Veto")
	KeyMaxBlockHeightPeriod    = []byte("MaxBlockHeightPeriod")
)

// ParamTable for gov module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&GovParams{})
}

// mint parameters
type GovParams struct {
	MaxDepositPeriod        time.Duration `json:"max_deposit_period"`          //  Maximum period for okb holders to deposit on a text proposal. Initial value: 2 days
	MinDeposit              sdk.DecCoins  `json:"min_deposit"`                 //  Minimum deposit for a critical text proposal to enter voting period.
	VotingPeriod            time.Duration `json:"voting_period"`               //  Length of the critical voting period for text proposal.
	DexListMaxDepositPeriod time.Duration `json:"dex_list_max_deposit_period"` //  Maximum period for okb holders to deposit on a dex list proposal. Initial value: 2 days
	DexListMinDeposit       sdk.DecCoins  `json:"dex_list_min_deposit"`        //  Minimum deposit for a critical dex list proposal to enter voting period.
	DexListVotingPeriod     time.Duration `json:"dex_list_voting_period"`      //  Length of the critical voting period for dex list proposal.
	DexListVoteFee          sdk.DecCoins  `json:"dex_list_vote_fee"`           //  Fee used for voting dex list proposal
	DexListMaxBlockHeight   uint64        `json:"dex_list_max_block_height"`   //  block height for dex list can not be greater than DexListMaxBlockHeight
	DexListFee              sdk.DecCoins  `json:"dex_list_fee"`                //  fee for dex list
	DexListExpireTime       time.Duration `json:"dex_list_expire_time"`        //  expire time for dex list
	Quorum                  sdk.Dec       `json:"quorum"`                      //
	Threshold               sdk.Dec       `json:"threshold"`                   //  Minimum proportion of Yes votes for proposal to pass. Initial value: 0.5
	Veto                    sdk.Dec       `json:"veto"`                        //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
	MaxBlockHeightPeriod    int64         `json:"max_block_height_period"`
}

func (p *GovParams) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyMaxDepositPeriod, &p.MaxDepositPeriod},
		{KeyMinDeposit, &p.MinDeposit},
		{KeyVotingPeriod, &p.VotingPeriod},
		{KeyDexListMaxDepositPeriod, &p.DexListMaxDepositPeriod},
		{KeyDexListMinDeposit, &p.DexListMinDeposit},
		{KeyDexListVotingPeriod, &p.DexListVotingPeriod},
		{KeyQuorum, &p.Quorum},
		{KeyThreshold, &p.Threshold},
		{KeyVeto, &p.Veto},
		{KeyMaxBlockHeightPeriod, &p.MaxBlockHeightPeriod},
		{KeyDexListVoteFee, &p.DexListVoteFee},
		{KeyDexListMaxBlockHeight, &p.DexListMaxBlockHeight},
		{KeyDexListFee, &p.DexListFee},
		{KeyDexListExpireTime, &p.DexListExpireTime},

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
		minDeposit, err := ParseDecCoins(value)
		fmt.Println(err)
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
	case string(KeyDexListMaxDepositPeriod):
		dexListMaxDepositPeriod, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateComplaintRetrospect(complaintRetrospect); err != nil {
		// 	return nil, err
		// }
		return dexListMaxDepositPeriod, nil
	case string(KeyDexListMinDeposit):
		dexListMinDeposit, err := ParseDecCoins(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateArbitrationTimeLimit(arbitrationTimeLimit); err != nil {
		// 	return nil, err
		// }
		return dexListMinDeposit, nil
	case string(KeyDexListVotingPeriod):
		dexListVotingPeriod, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateComplaintRetrospect(complaintRetrospect); err != nil {
		// 	return nil, err
		// }
		return dexListVotingPeriod, nil
	case string(KeyDexListVoteFee):
		fee, err := ParseDecCoins(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return fee, nil
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
	case string(KeyMaxBlockHeightPeriod):
		maxBlockHeightPeriod, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		return maxBlockHeightPeriod, nil
	case string(KeyDexListMaxBlockHeight):
		maxBlockHeight, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		return maxBlockHeight, nil
	case string(KeyDexListFee):
		fee, err := ParseDecCoins(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		return fee, nil
	case string(KeyDexListExpireTime):
		dexListExpireTime, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		return dexListExpireTime, nil
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
	case string(KeyDexListMaxDepositPeriod):
		err := cdc.UnmarshalJSON(bytes, &p.DexListMaxDepositPeriod)
		return p.DexListMaxDepositPeriod.String(), err
	case string(KeyDexListMinDeposit):
		err := cdc.UnmarshalJSON(bytes, &p.DexListMinDeposit)
		return p.DexListMinDeposit.String(), err
	case string(KeyDexListVotingPeriod):
		err := cdc.UnmarshalJSON(bytes, &p.DexListVotingPeriod)
		return p.DexListVotingPeriod.String(), err
	case string(KeyDexListVoteFee):
		err := cdc.UnmarshalJSON(bytes, &p.DexListVoteFee)
		return p.DexListVoteFee.String(), err
	case string(KeyQuorum):
		err := cdc.UnmarshalJSON(bytes, &p.Quorum)
		return p.Quorum.String(), err
	case string(KeyThreshold):
		err := cdc.UnmarshalJSON(bytes, &p.Threshold)
		return p.Threshold.String(), err
	case string(KeyVeto):
		err := cdc.UnmarshalJSON(bytes, &p.Veto)
		return p.Veto.String(), err
	case string(KeyMaxBlockHeightPeriod):
		err := cdc.UnmarshalJSON(bytes, &p.MaxBlockHeightPeriod)
		return strconv.FormatInt(p.MaxBlockHeightPeriod, 10), err
	case string(KeyDexListMaxBlockHeight):
		err := cdc.UnmarshalJSON(bytes, &p.DexListMaxBlockHeight)
		return strconv.FormatUint(p.DexListMaxBlockHeight, 10), err
	case string(KeyDexListFee):
		err := cdc.UnmarshalJSON(bytes, &p.DexListFee)
		return p.DexListFee.String(), err
	case string(KeyDexListExpireTime):
		err := cdc.UnmarshalJSON(bytes, &p.DexListExpireTime)
		return p.DexListExpireTime.String(), err
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

func (p GovParams) String() string {
	return fmt.Sprintf(`Deposit Params:
	Min Deposit:        %s
	Max Deposit Period: %s`, p.MinDeposit, p.MaxDepositPeriod) + "\n" +
		fmt.Sprintf(`DexList Deposit Params:
	DexList Min Deposit:        %s
	DexList Max Deposit Period: %s`, p.DexListMinDeposit, p.DexListMaxDepositPeriod) + "\n" +
		fmt.Sprintf(`DexList Voting Params:
	DexList Voting Period:      %s
	DexList Voting Fee:         %s`, p.DexListVotingPeriod, p.DexListVoteFee.String())  + "\n" +
		fmt.Sprintf(`DexList Params:
	DexList Max BlockHeight     %d
	DexList Fee                 %s
	DexList Expire Time         %s`,  p.DexListMaxBlockHeight, p.DexListFee.String(), p.DexListExpireTime.String()) + "\n" +
		fmt.Sprintf(`Tally Params:
	Quorum:             %s
	Threshold:          %s
	Veto:               %s`, p.Quorum, p.Threshold, p.Veto) + "\n" +
		fmt.Sprintf(`Voting Params:
	Voting Period:      %s`, p.VotingPeriod) + "\n" +
		fmt.Sprintf("MaxBlockHeightPeriod:    %d", p.MaxBlockHeightPeriod)
}

func NewGovParams(vp VotingParams, tp TallyParams, dp DepositParams, maxBlockHeightPeriod int64, dldp DexListDepositParams, dlvp DexListVotingParams, dlp DexListParams) GovParams {
	return GovParams{
		MaxDepositPeriod:        dp.MaxDepositPeriod,
		MinDeposit:              dp.MinDeposit,
		VotingPeriod:            vp.VotingPeriod,
		DexListMaxDepositPeriod: dldp.MaxDepositPeriod,
		DexListMinDeposit:       dldp.MinDeposit,
		DexListVotingPeriod:     dlvp.VotingPeriod,
		DexListVoteFee:          dlvp.VotingFee,
		DexListMaxBlockHeight:   dlp.MaxBlockHeight,
		DexListFee:              dlp.Fee,
		DexListExpireTime:       dlp.ExpireTime,
		Quorum:                  tp.Quorum,
		Threshold:               tp.Threshold,
		Veto:                    tp.Veto,
		MaxBlockHeightPeriod:    maxBlockHeightPeriod,
	}
}

// default minting module parameters
func DefaultParams() GovParams {
	//minDepositTokens := token.ToUnit(1000)
	var minDeposit = sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))}
	var dexListMinDeposit = sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(20000))}
	var votingFee = sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0))}
	var dexListFee = sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000))}

	return GovParams{
		MaxDepositPeriod:        DefaultPeriod,
		MinDeposit:              sdk.NewDecCoins(minDeposit),
		VotingPeriod:            DefaultPeriod,
		DexListMaxDepositPeriod: DefaultPeriod,
		DexListMinDeposit:       sdk.NewDecCoins(dexListMinDeposit),
		DexListVotingPeriod:     DefaultPeriod,
		DexListVoteFee:          sdk.NewDecCoins(votingFee),
		DexListMaxBlockHeight:   10000,
		DexListFee:              sdk.NewDecCoins(dexListFee),
		DexListExpireTime:       time.Hour * 24,
		Quorum:                  sdk.NewDecWithPrec(334, 3),
		Threshold:               sdk.NewDecWithPrec(5, 1),
		Veto:                    sdk.NewDecWithPrec(334, 3),
		MaxBlockHeightPeriod:    100000,
	}
}

func validateParams(p GovParams) sdk.Error {
	if err := validateDepositParams(DepositParams{
		MaxDepositPeriod: p.MaxDepositPeriod,
		MinDeposit:       p.MinDeposit,
	}); err != nil {
		return err
	}

	if err := validateDexListDepositParams(DexListDepositParams{
		MaxDepositPeriod: p.DexListMaxDepositPeriod,
		MinDeposit:       p.DexListMinDeposit,
	}); err != nil {
		return err
	}

	if err := validatorVotingParams(VotingParams{
		VotingPeriod: p.VotingPeriod,
	}); err != nil {
		return err
	}

	if err := validatorDexListVotingParams(DexListVotingParams{
		VotingPeriod: p.DexListVotingPeriod,
		VotingFee:    p.DexListVoteFee,
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

	if p.MaxBlockHeightPeriod < 0 {
		return ErrInvalMaxHeightPeriod(DefaultCodespace)
	}

	if err := validateDexListParams(DexListParams{
		MaxBlockHeight:    p.DexListMaxBlockHeight,
		Fee: p.DexListFee,
		ExpireTime: p.DexListExpireTime,
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
	MinDeposit       sdk.DecCoins
	MaxDepositPeriod time.Duration
}

func (dp DepositParams) String() string {
	return fmt.Sprintf(`Deposit Params:
  Min Deposit:        %s
  Max Deposit Period: %s`, dp.MinDeposit, dp.MaxDepositPeriod)
}

type DexListDepositParams struct {
	MinDeposit       sdk.DecCoins
	MaxDepositPeriod time.Duration
}

func (dp DexListDepositParams) String() string {
	return fmt.Sprintf(`DexListDeposit Params:
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

type DexListVotingParams struct {
	VotingPeriod time.Duration `json:"voting_period"` //  Length of the voting period.
	VotingFee    sdk.DecCoins  `json:"voting_fee"`    //  Fee for voting dex list proposal
}

func (dlvp DexListVotingParams) String() string {
	return fmt.Sprintf(`DexListVoting Params:
  Voting Period:      %s
  Voting Fee:         %s`, dlvp.VotingPeriod, dlvp.VotingFee.String())
}

type DexListParams struct {
	MaxBlockHeight uint64        `json:"max_block_height"` //  Max block height can be set for dex list.
	Fee            sdk.DecCoins  `json:"fee"`              //  Fee for dex
	ExpireTime     time.Duration `json:""`
}

func (dlp DexListParams) String() string {
	return fmt.Sprintf(`DexList Params:
  Max BlockHeight:      %d
  Fee:                  %s
  ExpireTime            %s`, dlp.MaxBlockHeight, dlp.Fee.String(), dlp.ExpireTime.String())
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

func validateDexListDepositParams(dldp DexListDepositParams) sdk.Error {
	if !dldp.MinDeposit.IsValid() {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositDenom, fmt.Sprintf("Governance dex list deposit amount must be a valid sdk.Coins amount, is %s",
			dldp.MinDeposit.String()))
	}
	return nil
}

func validatorVotingParams(vp VotingParams) sdk.Error {
	if vp.VotingPeriod < 20*time.Second || vp.VotingPeriod > 7*24*time.Hour {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidVotingPeriod, fmt.Sprintf("VotingPeriod (%s) should be between 20s and 1 week", vp.VotingPeriod.String()))
	}
	return nil
}

func validatorDexListVotingParams(dlvp DexListVotingParams) sdk.Error {
	if dlvp.VotingPeriod < 20*time.Second || dlvp.VotingPeriod > 7*24*time.Hour {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidVotingPeriod, fmt.Sprintf("Dex List VotingPeriod (%s) should be between 20s and 1 week", dlvp.VotingPeriod.String()))
	}
	if len(dlvp.VotingFee) != 1 || dlvp.VotingFee[0].Denom != sdk.DefaultBondDenom {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidVotingPeriod, fmt.Sprintf("Dex List VotingFee (%s) should be only in okb", dlvp.VotingFee.String()))
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

func validateDexListParams(dlp DexListParams) sdk.Error {
	if len(dlp.Fee) != 1 || dlp.Fee[0].Denom != sdk.DefaultBondDenom {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidVotingPeriod, fmt.Sprintf("Dex List Fee (%s) should be only in okb", dlp.Fee.String()))
	}
	return nil
}
