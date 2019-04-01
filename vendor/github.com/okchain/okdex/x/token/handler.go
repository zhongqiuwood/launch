package token

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	//"encoding/json"
	"bytes"
	"sort"
)

const (
	Suffix = "_freeze"
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

		default:
			errMsg := fmt.Sprintf("Unrecognized token Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgTokenIssue(ctx sdk.Context, keeper Keeper, msg MsgTokenIssue) sdk.Result {
	token := Token{
		Name:        msg.Name,
		Symbol:      msg.Symbol,
		TotalSupply: msg.TotalSupply,
		Owner:       msg.Owner,
		Mintable:    msg.Mintable,
	}

	// check whether the token exists
	t := keeper.GetTokenInfo(ctx, msg.Symbol)
	if t.Symbol != "" {
		return sdk.ErrInternal("The token already exists ").Result()
	}

	coins := keeper.coinKeeper.GetCoins(ctx, msg.Owner)
	//newCoins := keeper.coinKeeper.GetCoins(ctx, msg.Owner)
	coins = append(coins, msg.Tokens...)
	sort.Sort(ByDenom(coins))

	err := keeper.coinKeeper.SetCoins(ctx, msg.Owner, coins)
	if err != nil {
		return sdk.ErrInternal("Token issue error").Result()
	}

	// check whether the token is okb
	if token.Symbol != "okb" {
		feeCoins := FloatToCoins(FeeIssue)
		keeper.feeCollectionKeeper.AddCollectedFees(ctx, feeCoins)
		_, _, err = keeper.coinKeeper.SubtractCoins(ctx, msg.Owner, feeCoins)
		if err != nil {
			return sdk.ErrInsufficientCoins("Owner does not have enough okbs").Result()
		}
	}
	keeper.NewToken(ctx, token)

	return sdk.Result{}
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
			if coin.Amount.LT(ToUnit(msg.Amount)) {
				return sdk.ErrInsufficientCoins("Owner has insufficient coins").Result()
			}
		}
	}

	// subtract the coin
	var subCoins sdk.Coins
	coin := sdk.NewCoin(msg.Symbol, ToUnit(msg.Amount))
	subCoins = append(subCoins, coin)

	_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Owner, subCoins) // If so, deduct the Bid amount from the sender
	if err != nil {
		//fmt.Println(err)
		return sdk.ErrInsufficientCoins("Owner does not have enough coins").Result()
	}

	// update total supply
	token.TotalSupply -= msg.Amount
	keeper.NewToken(ctx, token)

	// token burn fees
	feeCoins := FloatToCoins(FeeBurn)
	keeper.feeCollectionKeeper.AddCollectedFees(ctx, feeCoins)
	_, _, err = keeper.coinKeeper.SubtractCoins(ctx, msg.Owner, feeCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins("Owner does not have enough okbs").Result()
	}

	return sdk.Result{}
}

func handleMsgTokenFreeze(ctx sdk.Context, keeper Keeper, msg MsgTokenFreeze) sdk.Result {
	coin := sdk.NewCoin(msg.Symbol, ToUnit(msg.Amount))
	var freezeCoins sdk.Coins
	freezeCoins = append(freezeCoins, coin)

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
	feeCoins := FloatToCoins(FeeFreeze)
	keeper.feeCollectionKeeper.AddCollectedFees(ctx, feeCoins)
	_, _, err = keeper.coinKeeper.SubtractCoins(ctx, msg.Owner, feeCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins("Owner does not have enough okbs").Result()
	}

	return sdk.Result{}
}

func handleMsgTokenUnfreeze(ctx sdk.Context, keeper Keeper, msg MsgTokenUnfreeze) sdk.Result {
	coin := sdk.NewCoin(msg.Symbol, ToUnit(msg.Amount))
	var unfreezeCoins sdk.Coins
	unfreezeCoins = append(unfreezeCoins, coin)

	// update unfreeze token
	oldCoins := keeper.GetFreezeTokens(ctx, msg.Owner)
	if oldCoins == nil {
		return sdk.ErrInsufficientCoins("Owner does not have enough unfreeze coins").Result()
	}

	newCoins, isNegative := oldCoins.SafeSub(unfreezeCoins)
	if isNegative {
		return sdk.ErrInsufficientCoins("Owner does not have enough unfreeze coins").Result()
	}
	sort.Sort(newCoins)
	if newCoins != nil {
		keeper.FreezeToken(ctx, msg.Owner, newCoins)
	} else {
		emptyCoin := sdk.NewCoin(msg.Symbol, sdk.NewInt(0))
		var emptyCoins sdk.Coins
		emptyCoins = append(emptyCoins, emptyCoin)
		keeper.FreezeToken(ctx, msg.Owner, emptyCoins)
	}

	// update account
	_, _, err := keeper.coinKeeper.AddCoins(ctx, msg.Owner, unfreezeCoins)
	if err != nil {
		fmt.Println(err)
		return sdk.ErrInsufficientCoins("Owner does not have enough coins").Result()
	}

	// token unfreeze fees
	feeCoins := FloatToCoins(FeeUnfreeze)
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
	feeCoins := FloatToCoins(FeeMint)
	keeper.feeCollectionKeeper.AddCollectedFees(ctx, feeCoins)
	_, _, err = keeper.coinKeeper.SubtractCoins(ctx, msg.Owner, feeCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins("Owner does not have enough okbs").Result()
	}

	return sdk.Result{}
}
