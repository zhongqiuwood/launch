package backend

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ok-chain/okchain/x/common"
	"github.com/ok-chain/okchain/x/token"
	"github.com/tendermint/tendermint/libs/log"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	orderKeeper    OrderKeeper  // The reference to the OrderKeeper to get deals
	tokenKeeper    TokenKeeper  // The reference to the TokenKeeper to get fee details
	cdc            *codec.Codec // The wire codec for binary encoding/decoding.
	orm            *ORM
	dealChan       chan *Deal
	endBlockerChan chan *EndBlockEvent
	stopChan       chan struct{}
	dealFlushChan  chan struct{}
	maintainConf   *MaintainConf
	logger         log.Logger
	latestTicker   *map[string]*Ticker
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(orderKeeper OrderKeeper, tokenKeeper TokenKeeper, cdc *codec.Codec, logger log.Logger, mainConfPath string) Keeper {
	k := Keeper{
		orderKeeper: orderKeeper,
		tokenKeeper: tokenKeeper,
		cdc:         cdc,
		logger:      logger,
	}

	k.maintainConf, _ = LoadMaintainConf(mainConfPath, "maintain.json")
	if k.maintainConf == nil {
		k.maintainConf = GetDefaultMaintainConfig()
	}
	k.logger.Info(fmt.Sprintf("[backend] maintain.conf %+v", k.maintainConf))

	if k.maintainConf.EnableBackend {
		orm, err := NewORM(k.maintainConf.LogSQL, k.maintainConf.Sqlite3Path, &logger, "backend.db")
		if err == nil {
			k.orm = orm
			k.endBlockerChan = make(chan *EndBlockEvent, k.maintainConf.BufferSize)
			k.dealChan = make(chan *Deal, k.maintainConf.BufferSize)
			k.stopChan = make(chan struct{})
			k.dealFlushChan = make(chan struct{}, 10)
			k.latestTicker = &map[string]*Ticker{}
			go consumeDeals(k.dealChan, k.stopChan, k.dealFlushChan, k.maintainConf, k.orm, &k.logger, k.latestTicker)
			go k.ConsumeEndBlockerEvent()
		}
	}

	return k
}

func storeDealAndKLine(ctx sdk.Context, keeper Keeper, blockHeight, timestamp int64) error {
	matchResultMap := keeper.orderKeeper.GetBlockMatchResult(ctx, ctx.BlockHeight())
	for product, matchResult := range matchResultMap.ResultMap {
		for _, record := range matchResult.Deals {
			order := keeper.orderKeeper.GetOrder(ctx, record.OrderId)

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
			keeper.StoreDeal(ctx, deal)

		}
		// TODO: compute kline and store
	}

	keeper.dealFlushChan <- struct{}{}

	return nil
}

func (k Keeper) ConsumeEndBlockerEvent() error {
	if k.maintainConf.EnableBackend == false {
		return nil
	}

	defer func() {
		r := recover()
		if r != nil {
			debug.PrintStack()
		}
	}()

	for {
		select {
		case event := <-k.endBlockerChan:
			storeNewOrders(event.ctx, k, event.blockHeight)
			updateOrders(event.ctx, k, event.blockHeight)
			storeDealAndKLine(event.ctx, k, event.blockHeight, event.timestamp)
			storeFeeDetails(event.ctx, k, event.blockHeight)

		case <-k.stopChan:
			return nil
		}
	}
	return nil
}

func (k Keeper) StoreDeal(ctx sdk.Context, deal *Deal) {

	//TODO: store Deal to db
	if k.logger == nil {
		k.logger = ctx.Logger().With("/x/backend")
	}

	if k.maintainConf.EnableBackend {
		k.dealChan <- deal
	}
}

func (k Keeper) StoreMatch(ctx sdk.Context, match *Match) {
	//TODO: store match to db
}

func (k Keeper) StoreKLineMin(ctx sdk.Context, kline *KlineM1) {
	//TODO: store kline to db
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

func consumeDeals(dealCh chan *Deal, stopCh chan struct{}, flushCh chan struct{}, conf *MaintainConf, orm *ORM, log *log.Logger, latestTicker *map[string]*Ticker) error {

	if log != nil {
		(*log).Debug("[backend] consumeDeals go routine started")
	}

	if conf.EnableBackend == false {
		return nil
	}

	defer func() {
		r := recover()
		if r != nil {
			debug.PrintStack()
		}
	}()

	ticker := time.NewTicker(time.Second * 60)
	//notifyChan := make(chan struct{}, 5)
	go generateKline1M(stopCh, conf, orm, log)

	refreshTSCh := make(chan int64, 100)
	go generateTicker(refreshTSCh, stopCh, orm, log, latestTicker)

	f := func() error {

		deals := []Deal{}
		dCount := len(dealCh)

		if dCount > 0 {
			(*log).Debug(fmt.Sprintf(
				"[backend] start consumeDeals#%+v# \n", len(deals)))
		} else {
			return nil
		}

		for i := 0; i < dCount; i++ {
			d := <-dealCh
			deals = append(deals, *d)
		}

		cnt, err := orm.AddDeals(&deals)
		if err != nil {
			(*log).Error(fmt.Sprintf("[backend] Expect to insert %d deals, inserted Count %d, err: %+v", len(deals), cnt, err))
		} else {
			(*log).Debug(fmt.Sprintf("[backend] Expect to insert %d deals, inserted Count %d", len(deals), cnt))
			if cnt > 0 {
				//notifyChan <- struct{}{}
				refreshTSCh <- deals[len(deals)-1].Timestamp
			}
		}
		return err
	}

	for {
		select {
		case <-flushCh:
			f()

		case <-ticker.C:
			f()

		case <-stopCh:
			return nil
		}
	}
	return nil
}

func generateKline1M(stop chan struct{}, conf *MaintainConf, orm *ORM, log *log.Logger) error {

	if log != nil {
		(*log).Debug("[backend] generateKline1M go routine started")
	}

	defer func() {
		r := recover()
		if r != nil {
			debug.PrintStack()
		}
	}()

	if conf.EnableBackend == false {
		return nil
	}

	startTS, endTS := int64(0), time.Now().Unix()
	anchorEndTS, _, err := orm.Deal2Kline1min(startTS, endTS)
	if err != nil {
		(*log).Debug(fmt.Sprintf("[backend] error: %+v \n", err))
	}

	time.Sleep(time.Duration(int(60-time.Now().Second()) * int(time.Second)))
	ticker := time.NewTicker(time.Second * 60)

	klineNotifyChans := generateSyncKlineMXChans()
	for freq, ntfCh := range *klineNotifyChans {
		go generateKlinesMX(ntfCh, stop, freq, conf, orm, log)
	}

	work := func() {
		crrtTS := time.Now().Unix()
		(*log).Debug(fmt.Sprintf("[backend] entering generateKline1M startTS: %d, endTS: %d\n", anchorEndTS, crrtTS))
		anchorStart, _, err := orm.Deal2Kline1min(anchorEndTS, crrtTS)
		if err != nil {
			(*log).Debug(fmt.Sprintf("[backend] error: %+v \n", err))

		} else {
			anchorEndTS = anchorStart

			for _, ch := range *klineNotifyChans {
				ch <- struct{}{}
			}
		}
	}

	lock := new(sync.Mutex)
	for {
		lock.Lock()
		select {
		//case <-notifyChan:
		//	work()
		case <-ticker.C:
			work()
		case <-stop:
			lock.Unlock()
			break

		}
		lock.Unlock()
	}
}

func generateSyncKlineMXChans() *map[int]chan struct{} {
	notifyChans := map[int]chan struct{}{}
	kMap := GetAllKlineMap()

	for freq, _ := range *kMap {
		notifyCh := make(chan struct{}, 1)
		notifyChans[freq] = notifyCh
	}

	return &notifyChans
}

func generateKlinesMX(notifyChan chan struct{}, stop chan struct{}, refreshInterval int, conf *MaintainConf, orm *ORM, log *log.Logger) error {

	if log != nil {
		(*log).Debug(fmt.Sprintf("[backend] generateKlineMX-#%d# go routine started", refreshInterval))
	}

	defer func() {
		r := recover()
		if r != nil {
			debug.PrintStack()
		}
	}()

	if conf.EnableBackend == false {
		return nil
	}

	destKName := GetKlineTableNameByFreq(refreshInterval)
	destK, err := NewKlineFactory(destKName, nil)
	destIKline := destK.(IKline)

	startTS, endTS := int64(0), time.Now().Unix()
	anchorEndTS, _, err := orm.MergeKlineM1(startTS, endTS, destIKline)
	if err != nil {
		(*log).Debug(fmt.Sprintf("[backend] error: %+v \n", err))
	}

	time.Sleep(time.Duration(int(60-time.Now().Second()) * int(time.Second)))
	interval := time.Duration(destIKline.GetFreqInSecond() * int(time.Second))
	(*log).Debug(fmt.Sprintf("[backend] duaration: %+v IKline: %+v \n", interval, destIKline))
	ticker := time.NewTicker(interval)

	work := func() {
		crrtTS := time.Now().Unix()
		(*log).Debug(fmt.Sprintf("[backend] entering generateKlinesMX-#%d# startTS: %d, endTS: %d\n",
			destIKline.GetFreqInSecond(), anchorEndTS, crrtTS))

		anchorStart, _, err := orm.MergeKlineM1(anchorEndTS, crrtTS, destIKline)
		if err != nil {
			(*log).Debug(fmt.Sprintf("[backend] error: %+v \n", err))

		} else {
			anchorEndTS = anchorStart
		}
	}

	for {
		select {
		case <-notifyChan:
			if time.Now().Unix() < anchorEndTS+int64(destIKline.GetFreqInSecond()) {
				break
			} else {
				work()
			}

		case <-ticker.C:
			work()

		case <-stop:
			break

		}
	}
}

func generateTicker(km1Chan chan int64, stop chan struct{}, orm *ORM, log *log.Logger, latestTicker *map[string]*Ticker) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			debug.PrintStack()
		}
	}()

	(*log).Debug(fmt.Sprintf("[backend] generateTicker go routine started"))

	defer func() {
		r := recover()
		if r != nil {
			debug.PrintStack()
		}
	}()

	for {
		select {
		case <-km1Chan:
			time.Sleep(time.Second)
			(*log).Debug(fmt.Sprintf("[backend] entering generateTicker"))
			crrTS := time.Now().Unix()
			if tickerMap, e := orm.KlineM1ToTicker(crrTS-120, crrTS); e == nil {
				if tickerMap != nil && len(tickerMap) > 0 {
					*latestTicker = tickerMap
					(*log).Debug(fmt.Sprintf("returnTickerMap %+v", latestTicker))
				} else {
					(*log).Debug(fmt.Sprintf("No product's deal refresh in [%d, %d), latestTicker: %+v", crrTS-120, crrTS, *latestTicker))
				}
			} else {
				(*log).Error(fmt.Sprintf("generateTicker error %+v, latestTickers %+v, returnTickers: %+v", e, *latestTicker, tickerMap))
				return e
			}

		case <-stop:
			break
		}
	}

	return nil
}

func storeFeeDetails(ctx sdk.Context, keeper Keeper, blockHeight int64) error {
	feeDetails := keeper.tokenKeeper.GetFeeDetailList(ctx, blockHeight)
	if len(feeDetails) > 0 {
		cnt, err := keeper.orm.AddFeeDetails(&feeDetails)
		if err != nil {
			keeper.logger.Error(fmt.Sprintf("[backend] Expect to insert %d feeDetails, inserted Count %d, err: %+v", len(feeDetails), cnt, err))
		} else {
			keeper.logger.Debug(fmt.Sprintf("[backend] Expect to insert %d feeDetails, inserted Count %d", len(feeDetails), cnt))
		}
		return err
	}
	return nil
}

func storeNewOrders(ctx sdk.Context, keeper Keeper, blockHeight int64) error {
	orderNum := keeper.orderKeeper.GetBlockOrderNum(ctx, blockHeight)
	var orders []Order
	var index int64 = 1
	for ; index < orderNum; index++ {
		orderId := common.FormatOrderId(blockHeight, index)
		order := keeper.orderKeeper.GetOrder(ctx, orderId)
		if order != nil {
			orderDb := Order{
				TxHash:         order.TxHash,
				OrderId:        order.OrderId,
				Sender:         order.Sender.String(),
				Product:        order.Product,
				Side:           order.Side,
				Price:          order.Price.String(),
				Quantity:       order.Quantity.String(),
				Status:         order.Status,
				FilledAvgPrice: order.FilledAvgPrice.String(),
				RemainQuantity: order.RemainQuantity.String(),
				Timestamp:      order.Timestamp,
			}
			orders = append(orders, orderDb)
		}
	}
	if len(orders) > 0 {
		cnt, err := keeper.orm.AddOrders(&orders)
		if err != nil {
			keeper.logger.Error(fmt.Sprintf("[backend] Expect to insert %d orders, inserted Count %d, err: %+v", len(orders), cnt, err))
		} else {
			keeper.logger.Debug(fmt.Sprintf("[backend] Expect to insert %d orders, inserted Count %d", len(orders), cnt))
		}
		return err
	}
	return nil
}

func updateOrders(ctx sdk.Context, keeper Keeper, blockHeight int64) error {
	var orders []Order
	orderIds := keeper.orderKeeper.GetUpdatedOrderIds(ctx, blockHeight)
	for _, orderId := range orderIds {
		order := keeper.orderKeeper.GetOrder(ctx, orderId)
		if order != nil {
			orderDb := Order{
				TxHash:         order.TxHash,
				OrderId:        order.OrderId,
				Sender:         order.Sender.String(),
				Product:        order.Product,
				Side:           order.Side,
				Price:          order.Price.String(),
				Quantity:       order.Quantity.String(),
				Status:         order.Status,
				FilledAvgPrice: order.FilledAvgPrice.String(),
				RemainQuantity: order.RemainQuantity.String(),
				Timestamp:      order.Timestamp,
			}
			orders = append(orders, orderDb)
		}
	}
	if len(orders) > 0 {
		cnt, err := keeper.orm.UpdateOrders(&orders)
		if err != nil {
			keeper.logger.Error(fmt.Sprintf("[backend] Expect to update %d orders, updated Count %d, err: %+v", len(orders), cnt, err))
		} else {
			keeper.logger.Debug(fmt.Sprintf("[backend] Expect to update %d orders, updated Count %d", len(orders), cnt))
		}
		return err
	}
	return nil
}

func (k Keeper) GetCandles(product string, granularity, size int) (r interface{}, err error) {
	m := GetAllKlineMap()
	candleType := (*m)[granularity]
	if candleType == "" || len(candleType) == 0 {
		return nil, fmt.Errorf("No %s found.", candleType)
	}

	klines, err := NewKlinesFactory(candleType)
	if err != nil {
		return nil, err
	} else {
		err := k.orm.getLatestKlinesByProduct(product, size, time.Now().Unix(), klines)
		return klines, err
	}

	return nil, err
}
