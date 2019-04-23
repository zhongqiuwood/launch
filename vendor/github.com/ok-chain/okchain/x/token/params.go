package token

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultParamspace = "token"
)

// Parameter keys
var (
	KeyListAsset              = []byte("ListAsset")
	KeyIssueAsset             = []byte("IssueAsset")
	KeyMintAsset              = []byte("MintAsset")
	KeyBurnAsset              = []byte("BurnAsset")
	KeyTransfer               = []byte("Transfer")
	KeyFreezeAsset            = []byte("FreezeAsset")
	KeyUnfreezeAsset          = []byte("UnfreezeAsset")
	KeyListPeriod             = []byte("ListPeriod")
	KeyListProposalMinDeposit = []byte("ListProposalMinDeposit")
)

var _ params.ParamSet = &Params{}

// mint parameters
type Params struct {
	ListAsset              sdk.Dec       `json:"list_asset"`                // Initial Coin Offering,Initial value:100000OKB
	IssueAsset             sdk.Dec       `json:"issue_asset"`               // Issue token,Initial value:20000OKB
	MintAsset              sdk.Dec       `json:"mint_asset"`                // Mint token,Initial value:2000OKB
	BurnAsset              sdk.Dec       `json:"burn_asset"`                // Burn token,Initial value:10OKB
	Transfer               sdk.Dec       `json:"transfer"`                  // Transfer,Initial value:0.0125OKB
	FreezeAsset            sdk.Dec       `json:"freeze_asset"`              // Freeze,Initial value:0.1OKB
	UnfreezeAsset          sdk.Dec       `json:"unfreeze_asset"`            // Unfreeze,Initial value:0.1OKB
	ListPeriod             time.Duration `json:"list_period"`               // Initial Coin Offering window,Initial value:24hours
	ListProposalMinDeposit sdk.Dec       `json:"list_proposal_min_deposit"` // Initial Coin Offering Min Deposit,Initial value:20000OKB
}

// ParamKeyTable for auth module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyListAsset, &p.ListAsset},
		{KeyIssueAsset, &p.IssueAsset},
		{KeyMintAsset, &p.MintAsset},
		{KeyBurnAsset, &p.BurnAsset},
		{KeyTransfer, &p.Transfer},
		{KeyFreezeAsset, &p.FreezeAsset},
		{KeyUnfreezeAsset, &p.UnfreezeAsset},
		{KeyListPeriod, &p.ListPeriod},
		{KeyListProposalMinDeposit, &p.ListProposalMinDeposit},
	}
}

func (p *Params) ValidateKV(key string, value string) (interface{}, sdk.Error) {
	switch key {
	case string(KeyListAsset):
		v, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return v, nil
	case string(KeyIssueAsset):
		v, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return v, nil
	case string(KeyMintAsset):
		v, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return v, nil
	case string(KeyBurnAsset):
		v, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return v, nil
	case string(KeyTransfer):
		v, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return v, nil
	case string(KeyFreezeAsset):
		v, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return v, nil
	case string(KeyUnfreezeAsset):
		v, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return v, nil
	case string(KeyListPeriod):
		v, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return v, nil
	case string(KeyListProposalMinDeposit):
		v, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return v, nil
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is not found", key))
	}
}

// // Equal returns a boolean determining if two Params types are identical.
// func (p Params) Equal(p2 Params) bool {
// 	bz1 := msgCdc.MustMarshalBinaryLengthPrefixed(&p)
// 	bz2 := msgCdc.MustMarshalBinaryLengthPrefixed(&p2)
// 	return bytes.Equal(bz1, bz2)
// }

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		ListAsset:              sdk.NewDecWithPrec(100000, 0),
		IssueAsset:             sdk.NewDecWithPrec(20000, 0),
		MintAsset:              sdk.NewDecWithPrec(2000, 0),
		BurnAsset:              sdk.NewDecWithPrec(100, 0),
		Transfer:               sdk.NewDecWithPrec(125, 4),
		FreezeAsset:            sdk.NewDecWithPrec(1, 1),
		UnfreezeAsset:          sdk.NewDecWithPrec(1, 1),
		ListPeriod:             24 * 60 * 60 * time.Second,
		ListProposalMinDeposit: sdk.NewDecWithPrec(20000, 0),
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("ListAsset: %s\n", p.ListAsset))
	sb.WriteString(fmt.Sprintf("IssueAsset: %s\n", p.IssueAsset))
	sb.WriteString(fmt.Sprintf("MintAsset: %s\n", p.MintAsset))
	sb.WriteString(fmt.Sprintf("BurnAsset: %s\n", p.BurnAsset))
	sb.WriteString(fmt.Sprintf("Transfer: %s\n", p.Transfer))
	sb.WriteString(fmt.Sprintf("FreezeAsset: %s\n", p.FreezeAsset))
	sb.WriteString(fmt.Sprintf("UnfreezeAsset: %s\n", p.UnfreezeAsset))
	sb.WriteString(fmt.Sprintf("ListPeriod: %s\n", p.ListPeriod))
	sb.WriteString(fmt.Sprintf("ListProposalMinDeposit: %s\n", p.ListProposalMinDeposit))
	return sb.String()
}
