package token

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ok-chain/okchain/x/common"

	//"encoding/json"
	"bytes"
	"sort"
	"strconv"
)

const (
	Suffix   = "_freeze"
	BaseCoin = "okb"
	// 90 billion
	TotalSupplyUpperbound = int64(90000000000)
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgTokenIssue:
			return handleMsgTokenIssue(ctx, keeper, msg)
		case MsgTokenBurn:
			return handleMsgTokenBurn(ctx, keeper, msg)

		case MsgTokenFreeze:
			return handleMsgTokenFreeze(ctx, keeper, msg)

		case MsgTokenUnfreeze:
			return handleMsgTokenUnfreeze(ctx, keeper, msg)

		case MsgTokenMint:
			return handleMsgTokenMint(ctx, keeper, msg)

		case MsgMultiSend:
			return handleMsgMultiSend(ctx, keeper, msg)

		case MsgSend:
			return handleMsgSend(ctx, keeper, msg)

		case MsgTokenTransfer:
			return handleMsgTokenTransfer(ctx, keeper, msg)

		default:
			errMsg := fmt.Sprintf("Unrecognized token Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgTokenIssue(ctx sdk.Context, keeper Keeper, msg MsgTokenIssue) sdk.Result {
	token := Token{
		Name:           msg.Name,
		Symbol:         msg.Symbol,
		OriginalSymbol: msg.OriginalSymbol,
		TotalSupply:    msg.TotalSupply,
		Owner:          msg.Owner,
		Mintable:       msg.Mintable,
	}

	if msg.TotalSupply >= TotalSupplyUpperbound {
		return sdk.ErrInternal("The number of the token exceeds the upper limit").Result()
	}

	// generate a random symbol
	for {
		token.Symbol = AddTokenSuffix(msg.OriginalSymbol)

		if !ValidCoinName(token.Symbol) {
			return sdk.ErrInvalidCoins("coin name not valid").Result()
		}
		msg.Symbol = token.Symbol

		// check whether the token exists
		t := keeper.GetTokenInfo(ctx, msg.Symbol)

		// find a symbol which doesn't exits
		if t.Symbol == "" {
			break
		}
	}

	coins := keeper.coinKeeper.GetCoins(ctx, msg.Owner)
	newCoin := sdk.NewCoin(msg.Symbol, ToUnit(msg.TotalSupply))
	var coinsToAdd []sdk.Coin
	coinsToAdd = append(coinsToAdd, newCoin)
	//newCoins := keeper.coinKeeper.GetCoins(ctx, msg.Owner)
	coins = append(coins, coinsToAdd...)
	sort.Sort(ByDenom(coins))

	err := keeper.coinKeeper.SetCoins(ctx, msg.Owner, coins)
	if err != nil {
		return sdk.ErrInternal("Token issue error").Result()
	}

	// check whether the token is okb
	if token.Symbol != BaseCoin {
		feeCoins := FloatToCoins(BaseCoin, FeeIssue)
		keeper.feeCollectionKeeper.AddCollectedFees(ctx, feeCoins)
		_, _, err = keeper.coinKeeper.SubtractCoins(ctx, msg.Owner, feeCoins)
		if err != nil {
			return sdk.ErrInsufficientCoins("Owner does not have enough okbs").Result()
		}
	}
	keeper.NewToken(ctx, token)

	return sdk.Result{
		Tags: sdk.NewTags("symbol", msg.Symbol),
	}

	//result := sdk.Result{}
	//result.Tags.AppendTag("symbol", msg.Symbol)
	//return result
}

func handleMsgTokenBurn(ctx sdk.Context, keeper Keeper, msg MsgTokenBurn) sdk.Result {
	token := keeper.GetTokenInfo(ctx, msg.Symbol)

	// check owner
	if !bytes.Equal(token.Owner.Bytes(), msg.Owner.Bytes()) {
		return sdk.ErrUnauthorized("Not the token's owner").Result()
	}

	// check balance
	myCoins := keeper.coinKeeper.GetCoins(ctx, msg.Owner)

	for _, coin := range myCoins {
		if coin.Denom == msg.Symbol {
			if coin.Amount.LT(sdk.MustNewDecFromStr(msg.Amount).RoundInt()) {
				return sdk.ErrInsufficientCoins("Owner has insufficient coins").Result()
			}
		}
	}

	// subtract the coin
	//var subCoins sdk.Coins
	//coin := sdk.NewCoin(msg.Symbol, ToUnit(msg.Amount))
	//subCoins = append(subCoins, coin)

	subCoins := FloatToCoins(msg.Symbol, msg.Amount)

	_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Owner, subCoins) // If so, deduct the Bid amount from the sender
	if err != nil {
		//fmt.Println(err)
		return sdk.ErrInsufficientCoins("Owner does not have enough coins").Result()
	}

	// update total supply
	// total supply int64
	// Todo: modify totalsupply(int64)
	supplyFloat, err2 := strconv.ParseFloat(msg.Amount, 32)
	if err2 != nil {
		return sdk.ErrInvalidCoins("invalid coins").Result()
	}
	token.TotalSupply -= int64(supplyFloat)
	keeper.NewToken(ctx, token)

	// token burn fees
	feeCoins := FloatToCoins(BaseCoin, FeeBurn)
	keeper.feeCollectionKeeper.AddCollectedFees(ctx, feeCoins)
	_, _, err = keeper.coinKeeper.SubtractCoins(ctx, msg.Owner, feeCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins("Owner does not have enough okbs").Result()
	}

	return sdk.Result{}
}

func handleMsgTokenFreeze(ctx sdk.Context, keeper Keeper, msg MsgTokenFreeze) sdk.Result {
	freezeCoins := FloatToCoins(msg.Symbol, msg.Amount)

	// update account
	_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Owner, freezeCoins) // If so, deduct the Bid amount from the sender
	if err != nil {
		fmt.Println(err)
		return sdk.ErrInsufficientCoins("Owner does not have enough coins").Result()
	}

	// update freeze token
	var newCoins sdk.Coins
	oldCoins := keeper.GetFreezeTokens(ctx, msg.Owner)
	if oldCoins == nil {
		newCoins = freezeCoins
	} else {
		newCoins = oldCoins.Add(freezeCoins)
	}

	sort.Sort(newCoins)
	keeper.FreezeToken(ctx, msg.Owner, newCoins)

	// token freeze fees
	feeCoins := FloatToCoins(BaseCoin, FeeFreeze)
	keeper.feeCollectionKeeper.AddCollectedFees(ctx, feeCoins)
	_, _, err = keeper.coinKeeper.SubtractCoins(ctx, msg.Owner, feeCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins("Owner does not have enough okbs").Result()
	}

	return sdk.Result{}
}

func handleMsgTokenUnfreeze(ctx sdk.Context, keeper Keeper, msg MsgTokenUnfreeze) sdk.Result {
	unfreezeCoins := FloatToCoins(msg.Symbol, msg.Amount)

	// update unfreeze token
	oldCoins := keeper.GetFreezeTokens(ctx, msg.Owner)
	if oldCoins == nil {
		return sdk.ErrInsufficientCoins("Owner does not have freeze coins").Result()
	}

	newCoins, isNegative := oldCoins.SafeSub(unfreezeCoins)
	if isNegative {
		return sdk.ErrInsufficientCoins(oldCoins.String() + "//" + unfreezeCoins.String()).Result()
		//return sdk.ErrInsufficientCoins("Owner does not have enough unfreeze coins").Result()
	}
	sort.Sort(newCoins)
	if newCoins != nil {
		keeper.FreezeToken(ctx, msg.Owner, newCoins)
	} else {
		keeper.ClearFreezeToken(ctx, msg.Owner)
	}

	// update account
	_, _, err := keeper.coinKeeper.AddCoins(ctx, msg.Owner, unfreezeCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins("Owner does not have enough coins").Result()
	}

	// token unfreeze fees
	feeCoins := FloatToCoins(BaseCoin, FeeUnfreeze)
	keeper.feeCollectionKeeper.AddCollectedFees(ctx, feeCoins)
	_, _, err = keeper.coinKeeper.SubtractCoins(ctx, msg.Owner, feeCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins("Owner does not have enough okbs").Result()
	}

	return sdk.Result{}
}

func handleMsgTokenMint(ctx sdk.Context, keeper Keeper, msg MsgTokenMint) sdk.Result {
	token := keeper.GetTokenInfo(ctx, msg.Symbol)
	// check owner
	if !bytes.Equal(token.Owner.Bytes(), msg.Owner.Bytes()) {
		return sdk.ErrUnauthorized("Not the token's owner").Result()
	}

	// check whether token is mintable
	if !token.Mintable {
		return sdk.ErrUnauthorized("token can't be minted").Result()
	}

	// modify total supply
	token.TotalSupply += msg.Amount
	if token.TotalSupply >= TotalSupplyUpperbound {
		return sdk.ErrInternal("The number of the token exceeds the upper limit").Result()
	}
	store := ctx.KVStore(keeper.tokenStoreKey)
	store.Set([]byte(token.Symbol), keeper.cdc.MustMarshalBinaryBare(token))

	coin := sdk.NewCoin(msg.Symbol, ToUnit(msg.Amount))
	var mintCoins sdk.Coins
	mintCoins = append(mintCoins, coin)
	_, _, err := keeper.coinKeeper.AddCoins(ctx, msg.Owner, mintCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins("Owner does not have enough coins").Result()
	}
	// token mint fees
	feeCoins := FloatToCoins(BaseCoin, FeeMint)
	keeper.feeCollectionKeeper.AddCollectedFees(ctx, feeCoins)
	_, _, err = keeper.coinKeeper.SubtractCoins(ctx, msg.Owner, feeCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins("Owner does not have enough okbs").Result()
	}

	return sdk.Result{}
}

func handleMsgMultiSend(ctx sdk.Context, keeper Keeper, msg MsgMultiSend) sdk.Result {
	for _, transferUnit := range msg.Transfers {
		err := keeper.SendCoins(ctx, msg.From, transferUnit.To, transferUnit.Coins)
		if err != nil {
			return sdk.ErrInsufficientCoins("Owner does not have enough okbs").Result()
		}
	}
	return sdk.Result{}
}

func handleMsgSend(ctx sdk.Context, keeper Keeper, msg MsgSend) sdk.Result {
	err := keeper.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return sdk.ErrInsufficientCoins("Owner does not have enough okbs").Result()
	}

	// token send fees
	feeCoins := FloatToCoins(BaseCoin, FeeTransfer)
	keeper.feeCollectionKeeper.AddCollectedFees(ctx, feeCoins)
	keeper.AddFeeDetail(ctx, msg.FromAddress.String(), feeCoins.String(), common.FeeTypeTransfer)
	_, _, err = keeper.coinKeeper.SubtractCoins(ctx, msg.FromAddress, feeCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins("Owner does not have enough okbs").Result()
	}

	return sdk.Result{}
}

func handleMsgTokenTransfer(ctx sdk.Context, keeper Keeper, msg MsgTokenTransfer) sdk.Result {
	tokenInfo := keeper.GetTokenInfo(ctx, msg.Symbol)

	if !tokenInfo.Owner.Equals(msg.FromAddress) {
		return sdk.ErrUnauthorized("not the token's owner").Result()
	}

	tokenInfo.Owner = msg.ToAddress
	keeper.NewToken(ctx, tokenInfo)
	return sdk.Result{}
}
