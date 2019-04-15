package order

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultParamspace = "order"
)

const (
	DefaultFeeNewOrder        = "0"      // okb
	DefaultFeeCancel          = "0.01"   // Equivalent okb
	DefaultFeeCancelNative    = "0.002"  // okb
	DefaultFeeExpire          = "0.01"   // Equivalent okb
	DefaultFeeExpireNative    = "0.002"  // okb
	DefaultFeeRateTrade       = "0.001"  // percentage
	DefaultFeeRateTradeNative = "0.0004" // percentage
)

// Parameter keys
var (
	KeyNewOrder           = []byte("NewOrder")
	KeyCancel             = []byte("Cancel")
	KeyCancelNative       = []byte("CancelNative")
	KeyExpire             = []byte("Expire")
	KeyExpireNative       = []byte("ExpireNative")
	KeyTradeFeeRate       = []byte("TradeFeeRate")
	KeyTradeFeeRateNative = []byte("TradeFeeRateNative")
)

var _ params.ParamSet = &Params{}

// mint parameters
type Params struct {
	NewOrder           sdk.Dec `json:"new_order"`             // 创建订单,Initial value:0
	Cancel             sdk.Dec `json:"cancel"`                // 取消订单,Initial value: 与0.01OKB等值
	CancelNative       sdk.Dec `json:"cancel_native"`         // 取消订单,Initial value: 0.002OKB
	Expire             sdk.Dec `json:"expire"`                // 订单超时,Initial value: 与0.01OKB等值
	ExpireNative       sdk.Dec `json:"expire_native"`         // 订单超时,Initial value: 0.002OKB
	TradeFeeRate       sdk.Dec `json:"trade_fee_rate"`        // 非OKB支付交易手续费率,Initial value: 0.001
	TradeFeeRateNative sdk.Dec `json:"trade_fee_rate_native"` // OKB支付交易手续费率,Initial value: 0.0004
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
		{KeyNewOrder, &p.NewOrder},
		{KeyCancel, &p.Cancel},
		{KeyCancelNative, &p.CancelNative},
		{KeyExpire, &p.Expire},
		{KeyExpireNative, &p.ExpireNative},
		{KeyTradeFeeRate, &p.TradeFeeRate},
		{KeyTradeFeeRateNative, &p.TradeFeeRateNative},
	}
}

func (p *Params) ValidateKV(key string, value string) (interface{}, sdk.Error) {
	switch key {
	case string(KeyNewOrder):
		v, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return v, nil
	case string(KeyCancel):
		v, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return v, nil
	case string(KeyCancelNative):
		v, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return v, nil
	case string(KeyExpire):
		v, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return v, nil
	case string(KeyExpireNative):
		v, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return v, nil
	case string(KeyTradeFeeRate):
		v, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		// if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
		// 	return nil, err
		// }
		return v, nil
	case string(KeyTradeFeeRateNative):
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
		NewOrder:           sdk.MustNewDecFromStr(DefaultFeeNewOrder),
		Cancel:             sdk.MustNewDecFromStr(DefaultFeeCancel),
		CancelNative:       sdk.MustNewDecFromStr(DefaultFeeCancelNative),
		Expire:             sdk.MustNewDecFromStr(DefaultFeeExpire),
		ExpireNative:       sdk.MustNewDecFromStr(DefaultFeeExpireNative),
		TradeFeeRate:       sdk.MustNewDecFromStr(DefaultFeeRateTrade),
		TradeFeeRateNative: sdk.MustNewDecFromStr(DefaultFeeRateTradeNative),
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("NewOrder: %s\n", p.NewOrder))
	sb.WriteString(fmt.Sprintf("Cancel: %s\n", p.Cancel))
	sb.WriteString(fmt.Sprintf("CancelNative: %s\n", p.CancelNative))
	sb.WriteString(fmt.Sprintf("Expire: %s\n", p.Expire))
	sb.WriteString(fmt.Sprintf("ExpireNative: %s\n", p.ExpireNative))
	sb.WriteString(fmt.Sprintf("TradeFeeRate: %s\n", p.TradeFeeRate))
	sb.WriteString(fmt.Sprintf("TradeFeeRateNative: %s\n", p.TradeFeeRateNative))
	return sb.String()
}
