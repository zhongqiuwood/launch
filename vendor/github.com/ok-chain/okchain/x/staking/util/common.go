package util

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// Denominations can be 3 ~ 16 characters long.
	reDnmString = `[a-z][a-z0-9]{2,15}`
	reDecAmt    = `[[:digit:]]*[\.]?[[:digit:]]+`
	reSpc       = `[[:space:]]*`
	reCoin      = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reDecAmt, reSpc, reDnmString))
)

// ParseCoin parses a cli input for one coin type, returning errors if invalid.
// This returns an error on an empty string as well.
func ParseCoin(coinStr string) (coin sdk.Coin, err error) {
	coinStr = strings.TrimSpace(coinStr)

	matches := reCoin.FindStringSubmatch(coinStr)
	if matches == nil {
		return sdk.Coin{}, fmt.Errorf("invalid coin expression: %s", coinStr)
	}

	denomStr, amountStr := matches[2], matches[1]

	amount, ok := NewIntFromStr(amountStr)
	if ok != nil {
		return sdk.Coin{}, fmt.Errorf("failed to parse coin amount: %s", amountStr)
	}

	if err := validateDenom(denomStr); err != nil {
		return sdk.Coin{}, fmt.Errorf("invalid denom cannot be used: %s", err)
	}

	return sdk.NewCoin(denomStr, amount), nil
}

func validateDenom(denom string) error {
	if strings.Compare(strings.ToLower(denom), "okb") != 0 {
		return errors.New("staking denom should be okb")
	}
	return nil
}

func NewIntFromStr(str string) (d sdk.Int, err sdk.Error) {
	if len(str) == 0 {
		return d, sdk.ErrUnknownRequest("decimal string is empty")
	}

	// first extract any negative symbol
	neg := false
	if str[0] == '-' {
		neg = true
		str = str[1:]
	}

	if len(str) == 0 {
		return d, sdk.ErrUnknownRequest("decimal string is empty")
	}

	strs := strings.Split(str, ".")
	lenDecs := 0
	combinedStr := strs[0]

	if len(strs) == 2 { // has a decimal place
		lenDecs = len(strs[1])
		if lenDecs == 0 || len(combinedStr) == 0 {
			return d, sdk.ErrUnknownRequest("bad decimal length")
		}
		combinedStr = combinedStr + strs[1]

	} else if len(strs) > 2 {
		return d, sdk.ErrUnknownRequest("too many periods to be a decimal string")
	}

	if lenDecs > sdk.Precision {
		return d, sdk.ErrUnknownRequest(
			fmt.Sprintf("too much precision, maximum %v, len decimal %v", sdk.Precision, lenDecs))
	}

	// add some extra zero's to correct to the Precision factor
	zerosToAdd := sdk.Precision - lenDecs
	zeros := fmt.Sprintf(`%0`+strconv.Itoa(zerosToAdd)+`s`, "")
	combinedStr = combinedStr + zeros

	combined, ok := new(big.Int).SetString(combinedStr, 10) // base 10
	if !ok {
		return d, sdk.ErrUnknownRequest(fmt.Sprintf("bad string to integer conversion, combinedStr: %v", combinedStr))
	}
	if neg {
		combined = new(big.Int).Neg(combined)
	}
	return sdk.NewIntFromBigInt(combined), nil
}

func HighPrecisionFromInt(input sdk.Int) sdk.Int {
	return sdk.NewIntFromBigInt(new(big.Int).Mul(input.BigInt(), new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))
}

func DecFromHighPrecision(input sdk.Int) sdk.Dec {
	return input.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))
}
