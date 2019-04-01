package backend

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	// The reference to the TokenKeeper to modify balances
	orderKeeper OrderKeeper
	// TODO: add db *gorm.DB

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(orderKeeper OrderKeeper, cdc *codec.Codec) Keeper {
	return Keeper{
		orderKeeper: orderKeeper,
		cdc:         cdc,
	}
}

func (k Keeper) StoreTrade(ctx sdk.Context, trade *Trade) {
	//TODO: store trade to db
}

func (k Keeper) StoreMatch(ctx sdk.Context, match *Match) {
	//TODO: store match to db
}

func (k Keeper) StoreKLineMin(ctx sdk.Context, kline *KLineMin) {
	//TODO: store kline to db
}

func (k Keeper) GetTrades(ctx sdk.Context, sender string) []Trade {
	return []Trade{{Sender: sender, OrderId: "test"}}
}
