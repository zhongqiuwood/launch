package gov

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// Denominations can be 3 ~ 16 characters long.
	reDnmString = `[a-z][a-z0-9]{2,15}`
	reAmt       = `[[:digit:]]+`
	reDecAmt    = `[[:digit:]]*\.?[[:digit:]]+`
	reSpc       = `[[:space:]]*`
	reDnm       = regexp.MustCompile(fmt.Sprintf(`^%s$`, reDnmString))
	reCoin      = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmt, reSpc, reDnmString))
	reDecCoin   = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reDecAmt, reSpc, reDnmString))
)

func NewCoinsFromDecCoins(decCoins sdk.DecCoins) sdk.Coins {
	cs := make(sdk.Coins, len(decCoins))
	for i, coin := range decCoins {
		cs[i] = NewCoinFromDecCoin(coin)
	}
	return cs
}

func NewCoinFromDecCoin(decCoin sdk.DecCoin) sdk.Coin {
	if decCoin.Amount.LT(sdk.NewDecFromBigInt(sdk.ZeroInt().BigInt())) {
		panic(fmt.Sprintf("negative decimal coin amount: %v\n", decCoin.Amount))
	}
	if strings.ToLower(decCoin.Denom) != decCoin.Denom {
		panic(fmt.Sprintf("denom cannot contain upper case characters: %s\n", decCoin.Denom))
	}

	return sdk.Coin{
		Denom:  decCoin.Denom,
		Amount: sdk.NewIntFromBigInt(decCoin.Amount.Int),
	}
}

func validateDenom(denom string) error {
	if !reDnm.MatchString(denom) {
		return errors.New("illegal characters")
	}
	return nil
}

func ParseDecCoins(coinsStr string) (coins sdk.DecCoins, err error) {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		return nil, nil
	}

	splitRe := regexp.MustCompile(",|;")
	coinStrs := splitRe.Split(coinsStr, -1)
	for _, coinStr := range coinStrs {
		coin, err := ParseDecCoin(coinStr)
		if err != nil {
			return nil, err
		}

		coins = append(coins, coin)
	}

	// sort coins for determinism
	coins.Sort()

	// validate coins before returning
	if !coins.IsValid() {
		return nil, fmt.Errorf("parsed decimal coins are invalid: %#v", coins)
	}

	return coins, nil
}

// ParseDecCoin parses a decimal coin from a string, returning an error if
// invalid. An empty string is considered invalid.
func ParseDecCoin(coinStr string) (coin sdk.DecCoin, err error) {
	coinStr = strings.TrimSpace(coinStr)

	matches := reDecCoin.FindStringSubmatch(coinStr)
	if matches == nil {
		return sdk.DecCoin{}, fmt.Errorf("invalid decimal coin expression: %s", coinStr)
	}

	amountStr, denomStr := matches[1], matches[2]

	amount, err := sdk.NewDecFromStr(amountStr)
	if err != nil {
		return sdk.DecCoin{}, errors.Wrap(err, fmt.Sprintf("failed to parse decimal coin amount: %s", amountStr))
	}

	if err := validateDenom(denomStr); err != nil {
		return sdk.DecCoin{}, fmt.Errorf("invalid denom cannot contain upper case characters or spaces: %s", err)
	}

	return sdk.NewDecCoinFromDec(denomStr, amount), nil
}
