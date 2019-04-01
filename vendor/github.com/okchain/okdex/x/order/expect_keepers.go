package order

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// expected token keeper
type TokenKeeper interface {
	HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error

	LockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) error
	UnlockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) error
	BurnLockedCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) error
	ReceiveLockedCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) error
}
