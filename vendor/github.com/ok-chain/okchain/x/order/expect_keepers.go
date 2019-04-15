package order

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ok-chain/okchain/x/token"
)

// expected token keeper
type TokenKeeper interface {
	// Token balance
	HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error

	LockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) error
	UnlockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) error
	BurnLockedCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) error
	ReceiveLockedCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) error

	// TokenPair
	GetTokenPair(ctx sdk.Context, product string) *token.TokenPair

	// Fee detail
	AddFeeDetail(ctx sdk.Context, from, fee, feeType string)
}
