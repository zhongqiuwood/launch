package token

import (
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/pkg/errors"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper          bank.Keeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	// The reference to the Param Keeper to get and set Global Params
	paramsKeeper params.Keeper
	// The reference to the Paramstore to get and set gov specific params
	paramSpace        params.Subspace
	tokenStoreKey     sdk.StoreKey // Unexposed key to access name store from sdk.Context
	freezeStoreKey    sdk.StoreKey
	lockStoreKey      sdk.StoreKey
	tokenPairStoreKey sdk.StoreKey
	//ownersStoreKey sdk.StoreKey // Unexposed key to access owners store from sdk.Context
	//amountStoreKey sdk.StoreKey // Unexposed key to access prices store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the token Keeper
func NewKeeper(coinKeeper bank.Keeper, paramsKeeper params.Keeper, paramSpace params.Subspace, feeCollectionKeeper auth.FeeCollectionKeeper, tokenStoreKey, freezeStoreKey sdk.StoreKey, lockStoreKey sdk.StoreKey, tokenPairStoreKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper:          coinKeeper,
		paramsKeeper:        paramsKeeper,
		paramSpace:          paramSpace.WithKeyTable(ParamKeyTable()),
		feeCollectionKeeper: feeCollectionKeeper,
		tokenStoreKey:       tokenStoreKey,
		freezeStoreKey:      freezeStoreKey,
		lockStoreKey:        lockStoreKey,
		tokenPairStoreKey:   tokenPairStoreKey,
		//originalStoreKey: originalStoreKey,
		//ownersStoreKey: ownersStoreKey,
		//amountStoreKey: amountStoreKey,
		cdc: cdc,
	}
}

func (k Keeper) GetTokenInfo(ctx sdk.Context, symbol string) Token {
	var token Token
	store := ctx.KVStore(k.tokenStoreKey)
	bz := store.Get([]byte(symbol))
	if bz == nil {
		return token
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &token)
	return token
}

func (k Keeper) GetTokensInfo(ctx sdk.Context) []Token {
	var tokens []Token
	store := ctx.KVStore(k.tokenStoreKey)
	iter := store.Iterator(nil, nil)
	for iter.Valid() {
		var token Token
		tokenBytes := iter.Value()
		k.cdc.MustUnmarshalBinaryBare(tokenBytes, &token)
		tokens = append(tokens, token)
		iter.Next()
	}
	return tokens
}

func (k Keeper) NewToken(ctx sdk.Context, token Token) {
	store := ctx.KVStore(k.tokenStoreKey)
	store.Set([]byte(token.Symbol), k.cdc.MustMarshalBinaryBare(token))
}

func (k Keeper) FreezeToken(ctx sdk.Context, acc sdk.AccAddress, coins sdk.Coins) {
	store := ctx.KVStore(k.freezeStoreKey)
	store.Set(acc.Bytes(), k.cdc.MustMarshalBinaryBare(coins))
}

func (k Keeper) GetFreezeTokens(ctx sdk.Context, acc sdk.AccAddress) sdk.Coins {
	store := ctx.KVStore(k.freezeStoreKey)
	coinsBytes := store.Get(acc.Bytes())
	if coinsBytes == nil {
		return nil
	}
	var coins sdk.Coins
	k.cdc.MustUnmarshalBinaryBare(coinsBytes, &coins)
	return coins
}

func (k Keeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return k.coinKeeper.HasCoins(ctx, addr, amt)
}
func (k Keeper) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error {
	_, _, err := k.coinKeeper.SubtractCoins(ctx, addr, amt)
	return err
}

func (k Keeper) LockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) error {
	// update account
	_, _, err := k.coinKeeper.SubtractCoins(ctx, addr, coins) // If so, deduct the Bid amount from the sender
	if err != nil {
		//fmt.Println(err)
		return errors.New("Owner does not have enough coins")
	}

	// update lock coins
	var newCoins sdk.Coins
	var oldCoins sdk.Coins
	store := ctx.KVStore(k.lockStoreKey)
	coinsBytes := store.Get(addr.Bytes())
	if coinsBytes == nil {
		newCoins = coins
	} else {
		k.cdc.MustUnmarshalBinaryBare(coinsBytes, &oldCoins)
		newCoins = oldCoins.Add(coins)
	}
	store.Set(addr.Bytes(), k.cdc.MustMarshalBinaryBare(newCoins))

	return nil
}

func (k Keeper) UnlockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) error {
	// check lockCoins
	var newCoins sdk.Coins
	var oldCoins sdk.Coins
	store := ctx.KVStore(k.lockStoreKey)
	coinsBytes := store.Get(addr.Bytes())
	if coinsBytes == nil {
		return errors.New("Owner does not have enough unlock coins")
	} else {
		var isNegative bool
		k.cdc.MustUnmarshalBinaryBare(coinsBytes, &oldCoins)
		newCoins, isNegative = oldCoins.SafeSub(coins)
		if isNegative {
			return errors.New("Owner does not have enough unlock coins")
		}
	}

	sort.Sort(newCoins)
	if newCoins != nil {
		store.Set(addr.Bytes(), k.cdc.MustMarshalBinaryBare(newCoins))
	}
	// update account
	_, _, err := k.coinKeeper.AddCoins(ctx, addr, coins)
	if err != nil {
		return errors.New("Add coins to Owner failed")
	}

	return nil
}

func (k Keeper) BurnLockedCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) error {
	store := ctx.KVStore(k.lockStoreKey)
	coinsBytes := store.Get(addr.Bytes())
	var oldCoins sdk.Coins
	k.cdc.MustUnmarshalBinaryBare(coinsBytes, &oldCoins)

	newCoins, isNegative := oldCoins.SafeSub(coins)
	if isNegative {
		return errors.New("Owner does not have enough unlock coins")
	}
	sort.Sort(newCoins)
	if newCoins != nil {
		store.Set(addr.Bytes(), k.cdc.MustMarshalBinaryBare(newCoins))
	}
	//_, _, err := k.coinKeeper.SubtractCoins(ctx, addr ,coins)
	return nil
}

func (k Keeper) ReceiveLockedCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) error {
	k.coinKeeper.AddCoins(ctx, addr, coins)
	return nil
}

// SaveTokenPair save the token pair to db
// key is base:quote
func (k Keeper) SaveTokenPair(ctx sdk.Context, tokenPair TokenPair) error {
	key := tokenPair.BaseAssetSymbol + ":" + tokenPair.QuoteAssetSymbol
	store := ctx.KVStore(k.tokenPairStoreKey)
	store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(tokenPair))
	return nil
}

// DropTokenPair drop the token pair
func (k Keeper) DropTokenPair(ctx sdk.Context, tokenPair TokenPair) error {
	key := tokenPair.BaseAssetSymbol + ":" + tokenPair.QuoteAssetSymbol
	store := ctx.KVStore(k.tokenPairStoreKey)
	store.Delete([]byte(key))
	return nil
}

// GetTokenPairs return all the token pairs
func (k Keeper) GetTokenPairs(ctx sdk.Context) []TokenPair {
	var tokenPairs []TokenPair
	store := ctx.KVStore(k.tokenPairStoreKey)
	iter := store.Iterator(nil, nil)
	for iter.Valid() {
		var tokenPair TokenPair
		tokenPairBytes := iter.Value()
		k.cdc.MustUnmarshalBinaryBare(tokenPairBytes, &tokenPair)
		tokenPairs = append(tokenPairs, tokenPair)
		iter.Next()
	}
	return tokenPairs
}

// SetCoins sets the coins at the addr.
func (k Keeper) SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return k.coinKeeper.SetCoins(ctx, addr, amt)
}

// GetCoins returns the coins at the addr.
func (k Keeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return k.coinKeeper.GetCoins(ctx, addr)
}

// get inflation params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) (params Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// set inflation params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
