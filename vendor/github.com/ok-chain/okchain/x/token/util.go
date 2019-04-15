package token

import (
	"math/big"

	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto"
	"regexp"
	"strings"
	"time"
	"sort"
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

var (
	//reDnmString = `[a-z][a-z0-9]{2,15}`
	//reDnmString = `[a-z0-9]{3,6}(\.O)?(\-)[a-z0-9]{3}`
	reCoinString = `[a-z0-9]{3,6}(\.o)?(\-)[a-z0-9]{3}`
	regCoin      = regexp.MustCompile(fmt.Sprintf(`^%s$`, reCoinString))
)

type BaseAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         sdk.Coins      `json:"coins"`
	PubKey        crypto.PubKey  `json:"public_key"`
	AccountNumber uint64         `json:"account_number"`
	Sequence      uint64         `json:"sequence"`
}

type DecAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         sdk.DecCoins   `json:"coins"`
	PubKey        crypto.PubKey  `json:"public_key"`
	AccountNumber uint64         `json:"account_number"`
	Sequence      uint64         `json:"sequence"`
}

// String implements fmt.Stringer
func (acc DecAccount) String() string {
	var pubkey string

	if acc.PubKey != nil {
		pubkey = sdk.MustBech32ifyAccPub(acc.PubKey)
	}

	return fmt.Sprintf(`Account:
 Address:       %s
 Pubkey:        %s
 Coins:         %v
 AccountNumber: %d
 Sequence:      %d`,
		acc.Address, pubkey, acc.Coins, acc.AccountNumber, acc.Sequence,
	)
}

func ToUnit(amount int64) sdk.Int {
	return sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(amount), new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))
}

func ValidCoinName(name string) bool {
	if reDnm.MatchString(name) {
		return true
	}
	return false
}

func AddTokenSuffix(name string) string {
	// will open suffix before production release
	return name
	timestamp := time.Now().Unix()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(timestamp))
	nameSuffix := base32.StdEncoding.EncodeToString(b)
	return name + "-" + strings.ToLower(nameSuffix[0:3])

}

func FloatToCoins(denom, fee string) sdk.Coins {
	feeDec := sdk.MustNewDecFromStr(fee)
	coin := sdk.NewCoin(denom, sdk.NewIntFromBigInt(feeDec.Int))
	var coins sdk.Coins
	coins = append(coins, coin)
	return coins
}

// AmountToCoins
// amount: 1:BNB,2:BTC
func AmountToCoins(amount string) sdk.Coins {
	var res sdk.Coins
	coinStrs := strings.Split(amount, ",")
	for _, coinStr := range coinStrs {
		coin := strings.Split(coinStr, ":")
		if len(coin) == 2 {
			var c sdk.Coin
			c.Denom = coin[1]
			coinDec := sdk.MustNewDecFromStr(coin[0])
			c.Amount = sdk.NewIntFromBigInt(coinDec.Int)
			res = append(res, c)
		}
	}
	return res
}

// format: [{"to": "addr", "amount": "1:BNB,2:BTC"}, ...]
// to []TransferUnit
func StrToTransfers(str string) (transfers []TransferUnit, err error) {
	var transfer []Transfer
	err = json.Unmarshal([]byte(str), &transfer)
	if err != nil {
		return transfers, err
	}

	for _, trans := range transfer {
		var t TransferUnit
		to, err := sdk.AccAddressFromBech32(trans.To)
		if err != nil {
			return transfers, err
		}
		t.To = to
		t.Coins = AmountToCoins(trans.Amount)
		transfers = append(transfers, t)
	}
	return transfers, nil
}

func BaseAccountToDecAccount(account auth.BaseAccount) DecAccount {
	var decCoins sdk.DecCoins
	for _, coin := range account.Coins {
		dec := sdk.NewDecFromBigIntWithPrec(coin.Amount.BigInt(), sdk.Precision)
		decCoin := sdk.NewDecCoinFromDec(coin.Denom, dec)
		decCoins = append(decCoins, decCoin)
	}
	decAccount := DecAccount{
		Address:       account.Address,
		PubKey:        account.PubKey,
		Coins:         decCoins,
		AccountNumber: account.AccountNumber,
		Sequence:      account.Sequence,
	}
	return decAccount
}

func validateDenom(denom string) error {
	if !reDnm.MatchString(denom) {
		return errors.New("illegal characters")
	}
	return nil
}

// ParseCoins will parse out a list of coins separated by commas.
// If nothing is provided, it returns nil Coins.
// Returned coins are sorted.
func ParseCoins(coinsStr string) (coins sdk.Coins, err error) {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		return nil, nil
	}

	coinStrs := strings.Split(coinsStr, ",")
	for _, coinStr := range coinStrs {
		coin, err := ParseDecCoin(coinStr)
		if err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}

	// Sort coins for determinism.
	coins.Sort()

	// Validate coins before returning.
	if !coins.IsValid() {
		return nil, fmt.Errorf("parseCoins invalid: %#v", coins)
	}

	return coins, nil
}

// ParseCoin parses a cli input for one coin type, returning errors if invalid.
// This returns an error on an empty string as well.
func ParseDecCoin(coinStr string) (coin sdk.Coin, err error) {
	coinStr = strings.TrimSpace(coinStr)

	matches := reDecCoin.FindStringSubmatch(coinStr)
	if matches == nil {
		return sdk.Coin{}, fmt.Errorf("invalid coin expression: %s", coinStr)
	}

	denomStr, amountStr := matches[2], matches[1]

	//amount, ok := sdk.NewIntFromString(amountStr)
	amount, err := sdk.NewDecFromStr(amountStr)
	if err != nil {
		return sdk.Coin{}, fmt.Errorf("failed to parse coin amount %s: %s", amountStr, err.Error())
	}

	if err := validateDenom(denomStr); err != nil {
		return sdk.Coin{}, fmt.Errorf("invalid denom cannot contain upper case characters or spaces: %s", err)
	}

	coin = sdk.NewCoin(denomStr, sdk.NewIntFromBigInt(amount.Int))

	return coin, nil
}

func MergeCoinInfo(availableCoins, freezeCoins, lockCoins sdk.Coins) (coinsInfo CoinsInfo) {
	m := make(map[string]CoinInfo)

	for _, availableCoin := range availableCoins {
		coinInfo, ok := m[availableCoin.Denom]
		if ok {
			dec := sdk.NewDecFromBigIntWithPrec(availableCoin.Amount.BigInt(), sdk.Precision)
			coinInfo.Available = dec.String()
			m[availableCoin.Denom]  = coinInfo
		} else {
			coinInfo.Symbol = availableCoin.Denom
			dec := sdk.NewDecFromBigIntWithPrec(availableCoin.Amount.BigInt(), sdk.Precision)
			coinInfo.Available = dec.String()
			coinInfo.Freeze = "0"
			coinInfo.Locked = "0"
			m[availableCoin.Denom]  = coinInfo
		}
	}

	for _, freezeCoin := range freezeCoins {
		coinInfo, ok := m[freezeCoin.Denom]
		if ok {
			dec := sdk.NewDecFromBigIntWithPrec(freezeCoin.Amount.BigInt(), sdk.Precision)
			coinInfo.Freeze = dec.String()
			m[freezeCoin.Denom]  = coinInfo
		} else {
			coinInfo.Symbol = freezeCoin.Denom
			dec := sdk.NewDecFromBigIntWithPrec(freezeCoin.Amount.BigInt(), sdk.Precision)
			coinInfo.Freeze = dec.String()
			coinInfo.Available = "0"
			coinInfo.Locked = "0"
			m[freezeCoin.Denom]  = coinInfo
		}
	}

	for _, lockCoin := range lockCoins {
		coinInfo, ok := m[lockCoin.Denom]
		if ok {
			dec := sdk.NewDecFromBigIntWithPrec(lockCoin.Amount.BigInt(), sdk.Precision)
			coinInfo.Locked = dec.String()
			m[lockCoin.Denom]  = coinInfo
		} else {
			coinInfo.Symbol = lockCoin.Denom
			dec := sdk.NewDecFromBigIntWithPrec(lockCoin.Amount.BigInt(), sdk.Precision)
			coinInfo.Available = "0"
			coinInfo.Locked = dec.String()
			coinInfo.Freeze = "0"
			m[lockCoin.Denom]  = coinInfo
		}
	}

	for _, coinInfo := range m {
		coinsInfo = append(coinsInfo, coinInfo)
	}
	sort.Sort(coinsInfo)
	return coinsInfo
}
