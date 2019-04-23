package backend

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ok-chain/okchain/x/token"
	"github.com/tendermint/tendermint/libs/log"
	"strconv"
	"time"
	"github.com/ok-chain/okchain/x/common"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	orderKeeper    OrderKeeper  // The reference to the OrderKeeper to get deals
	tokenKeeper    TokenKeeper  // The reference to the TokenKeeper to get fee details
	cdc            *codec.Codec // The wire codec for binary encoding/decoding.
	orm            *ORM
	dealChan       chan *Deal
	endBlockerChan chan *EndBlockEvent
	txEventChan    chan *TxEvent
	txChan         chan *Transaction
	stopChan       chan struct{}
	dealFlushChan  chan struct{}
	txFlushChan    chan struct{}
	maintainConf   *MaintainConf
	logger         log.Logger
	latestTicker   *map[string]*Ticker
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(orderKeeper OrderKeeper, tokenKeeper TokenKeeper, cdc *codec.Codec, logger log.Logger, mainConf *MaintainConf) Keeper {
	k := Keeper{
		orderKeeper: orderKeeper,
		tokenKeeper: tokenKeeper,
		cdc:         cdc,
		logger:      logger,
		maintainConf: mainConf,
	}

	k.logger.Info(fmt.Sprintf("[backend] maintain.conf %+v", k.maintainConf))

	if k.maintainConf.EnableBackend {
		orm, err := NewORM(k.maintainConf.LogSQL, k.maintainConf.Sqlite3Path, "backend.db", &logger)
		if err == nil {
			k.orm = orm
			k.endBlockerChan = make(chan *EndBlockEvent, k.maintainConf.BufferSize)
			k.txEventChan = make(chan *TxEvent, k.maintainConf.BufferSize)
			k.txChan = make(chan *Transaction, k.maintainConf.BufferSize)
			k.txFlushChan = make(chan struct{}, 10)
			k.dealChan = make(chan *Deal, k.maintainConf.BufferSize)
			k.stopChan = make(chan struct{})
			k.dealFlushChan = make(chan struct{}, 10)
			k.latestTicker = &map[string]*Ticker{}

			go k.ConsumeDeal()
			go k.ConsumeEndBlockerEvent()
			go k.ConsumeTxEvent()
			go k.ConsumeTx()
		}
	}

	return k
}

func (k *Keeper) Stop() {
	defer common.PrintStackIfPainic()
	if k.stopChan != nil {
		close(k.stopChan)
	}
	if k.orm != nil {
		k.orm.Close()
	}
}

func storeDealAndKLine(ctx sdk.Context, keeper Keeper, blockHeight, timestamp int64) error {
	result := keeper.orderKeeper.GetBlockMatchResult(ctx, ctx.BlockHeight())
	if result == nil {
		return nil
	}
	for product, matchResult := range result.ResultMap {
		for _, record := range matchResult.Deals {
			order := keeper.orderKeeper.GetOrder(ctx, record.OrderId)
			keeper.logger.Debug(fmt.Sprintf("[backend] storeDealAndKLine %+v", order))

			price, _ := strconv.ParseFloat(matchResult.Price.String(), 64)
			quantity, _ := strconv.ParseFloat(record.Quantity.String(), 64)

			deal := &Deal{
				BlockHeight: blockHeight,
				OrderId:     record.OrderId,
				Side:        record.Side,
				Sender:      order.Sender.String(),
				Product:     product,
				Price:       price,
				Quantity:    quantity,
				Timestamp:   timestamp,
			}
			keeper.ProduceDeal(deal)
		}
	}
	keeper.dealFlushChan <- struct{}{}
	return nil
}

func (k Keeper) GetDeals(ctx sdk.Context, sender, product string, offset, limit int) (*[]Deal, int) {
	return k.orm.GetDeals(sender, product, offset, limit)
}

func (k Keeper) GetFeeDetails(ctx sdk.Context, addr string, offset, limit int) (*[]token.FeeDetail, int) {
	return k.orm.GetFeeDetails(addr, offset, limit)
}

func (k Keeper) GetOrderList(ctx sdk.Context, addr, product string, open bool, offset, limit int) (*[]Order, int) {
	return k.orm.GetOrderList(addr, product, open, offset, limit)
}

func (k Keeper) GetTransactionList(ctx sdk.Context, addr string, txType, startTime, endTime int64, offset, limit int) (*[]Transaction, int) {
	return k.orm.GetTransactionList(addr, txType, startTime, endTime, offset, limit)
}

func (k Keeper) GetCandles(product string, granularity, size int) (r [][]string, err error) {
	if k.maintainConf.EnableBackend == false {
		return nil, fmt.Errorf("backend is not enabled, no candle found, maintian.conf: %+v", k.maintainConf)
	}

	m := GetAllKlineMap()
	candleType := (*m)[granularity]
	if candleType == "" || len(candleType) == 0 || (size < 0 || size > 1000) {
		return nil, fmt.Errorf("parameter's not correct, size: %d, granularity: %d", size, granularity)
	}

	klines, err := NewKlinesFactory(candleType)
	if err == nil {
		err := k.orm.getLatestKlinesByProduct(product, size, time.Now().Unix(), klines)
		iklines := ToIKlinesArray(klines, time.Now().Unix(), true)
		restData := ToRestfulData(&iklines, size)
		return restData, err
	}

	return nil, err
}

func (k Keeper) GetTickers(products []string, count int) []Ticker {
	tickers := []Ticker{}
	if k.latestTicker != nil {

		if len(products) > 0 {
			for _, p := range products {
				t := (*k.latestTicker)[p]
				if t != nil {
					tickers = append(tickers, *t)
				}
			}
		} else {
			for _, ticker := range *k.latestTicker {
				tickers = append(tickers, *ticker)
			}
		}
	}

	maxUpper := count
	if len(tickers) > 0 {
		if len(tickers) < maxUpper {
			maxUpper = len(tickers)
		}
		return tickers[0:maxUpper]
	} else {
		return tickers
	}
}
