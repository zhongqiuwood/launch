package backend

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/ok-chain/okchain/x/common"
	orderm "github.com/ok-chain/okchain/x/order"
	"github.com/ok-chain/okchain/x/token"
	"github.com/tendermint/tendermint/libs/log"
		"time"
)

// Produce data for txEventChan
func (k Keeper) ProduceTxEvent(ctx sdk.Context, tx *auth.StdTx, txHash string, timestamp int64) {
	if k.maintainConf.EnableBackend {
		k.txEventChan <- &TxEvent{ctx, tx, txHash, timestamp}
	}
}

// Consume data from txEventChan
func (k Keeper) ConsumeTxEvent() {

	defer common.PrintStackIfPainic()

	for {
		select {
		case txEvent := <-k.txEventChan:
			generateTx(txEvent.tx, txEvent.txHash, txEvent.ctx, k, txEvent.timestamp)

		case <-k.stopChan:
			return
		}
	}
}

// Produce data for txChan
func (k Keeper) ProduceTx(tx *Transaction) {
	if k.maintainConf.EnableBackend {
		k.txChan <- tx
	}
}

// Consume data from txChan
func (k Keeper) ConsumeTx() {
	defer common.PrintStackIfPainic()

	f := func() error {
		txs := []Transaction{}
		txChanLen := len(k.txChan)

		if txChanLen > 0 {
			k.logger.Debug(fmt.Sprintf(
				"[backend] start consumeTxs#%+v# \n", len(txs)))
		} else {
			return nil
		}

		for i := 0; i < txChanLen; i++ {
			tx := <-k.txChan
			txs = append(txs, *tx)
		}

		cnt, err := k.orm.AddTransactions(&txs)
		if err != nil {
			k.logger.Error(fmt.Sprintf("[backend] Expect to insert %d txs, inserted Count %d, err: %+v", len(txs), cnt, err))
		} else {
			k.logger.Debug(fmt.Sprintf("[backend] Expect to insert %d txs, inserted Count %d", len(txs), cnt))
		}
		return err
	}

	for {
		select {
		case <-k.txFlushChan:
			f()

		case <-k.stopChan:
			return
		}
	}
}

// Product data for endBlockerChan
func (k Keeper) ProductEndBlockerEvent(event *EndBlockEvent) {
	if k.maintainConf.EnableBackend {
		k.endBlockerChan <- event
	}
}

// Consume data from endBlockerChan
func (k Keeper) ConsumeEndBlockerEvent() {

	defer common.PrintStackIfPainic()

	for {
		select {
		case event := <-k.endBlockerChan:
			storeNewOrders(event.ctx, k, event.blockHeight)
			updateOrders(event.ctx, k, event.blockHeight)
			storeDealAndKLine(event.ctx, k, event.blockHeight, event.timestamp)
			storeFeeDetails(event.ctx, k, event.blockHeight)
			k.txFlushChan <- struct{}{}

		case <-k.stopChan:
			return
		}
	}
}

// Produce data for dealChan
func (k Keeper) ProduceDeal(deal *Deal) {
	if k.maintainConf.EnableBackend {
		k.dealChan <- deal
	}
}

// Consume data from dealChan
func (k Keeper) ConsumeDeal() {
	if k.logger != nil {
		k.logger.Debug("[backend] ConsumeDeal go routine started")
	}

	defer common.PrintStackIfPainic()

	ticker := time.NewTicker(time.Second * 60)

	ts := time.Now().Unix()
	UpdateTickersBuffer(k.orm, k.latestTicker, ts - 84400, ts)

	go generateKline1M(k.stopChan, k.maintainConf, k.orm, &k.logger)

	f := func() error {
		deals := []Deal{}
		dCount := len(k.dealChan)

		if dCount > 0 {
			k.logger.Debug(fmt.Sprintf(
				"[backend] start consumeDeals#%+v# \n", len(deals)))
		} else {
			return nil
		}

		for i := 0; i < dCount; i++ {
			d := <-k.dealChan
			deals = append(deals, *d)
		}

		cnt, err := k.orm.AddDeals(&deals)
		if err != nil {
			k.logger.Error(fmt.Sprintf("[backend] Expect to insert %d deals, inserted Count %d, err: %+v", len(deals), cnt, err))
		} else {
			k.logger.Debug(fmt.Sprintf("[backend] Expect to insert %d deals, inserted Count %d", len(deals), cnt))
			if cnt > 0 {
				ts := time.Now().Unix()
				UpdateTickersBuffer(k.orm, k.latestTicker, ts - 120, ts+1)
			}
		}
		return err
	}

	for {
		select {
		case <-k.dealFlushChan:
			f()

		case <-ticker.C:
			f()

		case <-k.stopChan:
			return
		}
	}
}

func buildTransactionsTransfer(msg token.MsgSend, txHash string, ctx sdk.Context, keeper Keeper, timestamp int64) (*Transaction, *Transaction) {
	decCoins := common.ConvertCoinsToDecCoins(msg.Amount)

	txFrom := &Transaction{
		TxHash:    txHash,
		Address:   msg.FromAddress.String(),
		Type:      TxTypeTransfer,
		Side:      TxSideFrom,
		Symbol:    decCoins[0].Denom,
		Quantity:  decCoins[0].Amount.String(),
		Fee:       sdk.DecCoin{Denom: common.ChainAsset, Amount: sdk.MustNewDecFromStr(token.FeeTransfer)}.String(), // TODO: get fee from params
		Timestamp: timestamp,
	}
	txTo := &Transaction{
		TxHash:    txHash,
		Address:   msg.ToAddress.String(),
		Type:      TxTypeTransfer,
		Side:      TxSideTo,
		Symbol:    decCoins[0].Denom,
		Quantity:  decCoins[0].Amount.String(),
		Fee:       sdk.DecCoin{Denom: common.ChainAsset, Amount: sdk.ZeroDec()}.String(),
		Timestamp: timestamp,
	}
	return txFrom, txTo
}

func buildTransactionNew(msg orderm.MsgNewOrder, txHash string, ctx sdk.Context, keeper Keeper, timestamp int64) *Transaction {
	side := TxSideBuy
	if msg.Side == "SELL" {
		side = TxSideSell
	}
	return &Transaction{
		TxHash:    txHash,
		Address:   msg.Sender.String(),
		Type:      TxTypeOrderNew,
		Side:      int64(side),
		Symbol:    msg.Product,
		Quantity:  msg.Quantity,
		Fee:       sdk.DecCoin{Denom: common.ChainAsset, Amount: sdk.ZeroDec()}.String(), // TODO: get fee from params
		Timestamp: timestamp,
	}
}

func buildTransactionCancel(msg orderm.MsgCancelOrder, txHash string, ctx sdk.Context, keeper Keeper, timestamp int64) *Transaction {
	order := keeper.orderKeeper.GetOrder(ctx, msg.OrderId)
	if order == nil {
		return nil
	}
	side := TxSideBuy
	if order.Side == "SELL" {
		side = TxSideSell
	}
	return &Transaction{
		TxHash:    txHash,
		Address:   order.Sender.String(),
		Type:      TxTypeOrderCancel,
		Side:      int64(side),
		Symbol:    order.Product,
		Quantity:  order.Quantity.String(),
		Fee:       order.GetExtraInfoWithKey(orderm.OrderExtraInfoKeyCancelFee),
		Timestamp: timestamp,
	}
}

func generateTx(tx *auth.StdTx, txHash string, ctx sdk.Context, keeper Keeper, timestamp int64) {
	for _, msg := range tx.GetMsgs() {
		switch msg.Type() {
		case "send": // token/send
			txFrom, txTo := buildTransactionsTransfer(msg.(token.MsgSend), txHash, ctx, keeper, timestamp)
			keeper.ProduceTx(txFrom)
			keeper.ProduceTx(txTo)
		case "new": // order/new
			transaction := buildTransactionNew(msg.(orderm.MsgNewOrder), txHash, ctx, keeper, timestamp)
			keeper.ProduceTx(transaction)
		case "cancel": // order/cancel
			transaction := buildTransactionCancel(msg.(orderm.MsgCancelOrder), txHash, ctx, keeper, timestamp)
			keeper.ProduceTx(transaction)
		default: // In other cases, do nothing
			continue
		}
	}
}


func generateKline1M(stop chan struct{}, conf *MaintainConf, orm *ORM, log *log.Logger) error {
	orm.Debug("[backend] generateKline1M go routine started")
	defer common.PrintStackIfPainic()

	startTS, endTS := int64(0), time.Now().Unix()
	anchorEndTS, _, err := orm.Deal2Kline1min(startTS, endTS)
	if err != nil {
		(*log).Debug(fmt.Sprintf("[backend] error: %+v \n", err))
	}

	time.Sleep(time.Duration(int(60-time.Now().Second()) * int(time.Second)))
	ticker := time.NewTicker(time.Second * 60)

	go cleanUpKlines(stop, orm, conf)

	klineNotifyChans := generateSyncKlineMXChans()
	for freq, ntfCh := range *klineNotifyChans {
		go generateKlinesMX(ntfCh, stop, freq, orm)
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

	for {
		select {
		case <-ticker.C:
			work()
		case <-stop:
			break

		}
	}
}

func generateSyncKlineMXChans() *map[int]chan struct{} {
	notifyChans := map[int]chan struct{}{}
	kMap := GetAllKlineMap()

	for freq, _ := range *kMap {
		if freq > 60 {
			notifyCh := make(chan struct{}, 1)
			notifyChans[freq] = notifyCh
		}
	}

	return &notifyChans
}

func generateKlinesMX(notifyChan chan struct{}, stop chan struct{}, refreshInterval int, orm *ORM) error {
	orm.Debug(fmt.Sprintf("[backend] generateKlineMX-#%d# go routine started", refreshInterval))


	destKName := GetKlineTableNameByFreq(refreshInterval)
	destK, err := NewKlineFactory(destKName, nil)
	destIKline := destK.(IKline)

	startTS, endTS := int64(0), time.Now().Unix()
	anchorEndTS, _, err := orm.MergeKlineM1(startTS, endTS, destIKline)
	if err != nil {
		orm.Debug(fmt.Sprintf("[backend] error: %s \n", err.Error()))
	}

	time.Sleep(time.Duration(int(60-time.Now().Second()) * int(time.Second)))
	interval := time.Duration(destIKline.GetFreqInSecond() * int(time.Second))
	orm.Debug(fmt.Sprintf("[backend] duaration: %+v IKline: %+v(%d s) \n", interval, destIKline, destIKline.GetFreqInSecond()))
	ticker := time.NewTicker(interval)

	work := func() {
		crrtTS := time.Now().Unix()
		orm.Debug(fmt.Sprintf("[backend] entering generateKlinesMX-#%d# startTS: %d, endTS: %d\n",
			destIKline.GetFreqInSecond(), anchorEndTS, crrtTS))

		anchorStart, _, err := orm.MergeKlineM1(anchorEndTS, crrtTS, destIKline)
		if err != nil {
			orm.Debug(fmt.Sprintf("[backend] error: %s \n", err.Error()))

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

func cleanUpKlines(stop chan struct{}, orm *ORM, conf *MaintainConf)  {
	orm.Debug(fmt.Sprintf("[backend] cleanUpKlines go routine started. MaintainConf: %+v", *conf))
	time.Sleep(time.Duration(int(60-time.Now().Second()) * int(time.Second)))
	interval := time.Duration(60 * int(time.Second))
	ticker := time.NewTicker(interval)

	work := func() {
		now := time.Now()
		strNow := now.Format("15:04:05")
		if strNow == conf.CleanUpsTime {

			m := GetAllKlineMap()
			for _, ktype := range *m {
				expiredDays := conf.CleanUpsKeptDays[ktype]
				if expiredDays != 0 {
					orm.Debug(fmt.Sprintf("[backend] entering cleanUpKlines, " +
						"fired time: %s(currentTS: %d), kline type: %s", conf.CleanUpsTime, now.Unix(), ktype))
					anchorTS := now.Add(-time.Duration(int(time.Second) * 1440 * expiredDays)).Unix()
					kline, _ := NewKlineFactory(ktype, nil)
					orm.DeleteKlineBefore(anchorTS, kline)
				}
			}
		}
	}

	for {
		select {
		case <-ticker.C:
			work()

		case <-stop:
			break

		}
	}
}

func UpdateTickersBuffer(orm *ORM, latestTicker *map[string]*Ticker, startTS, endTS int64) (err error) {
	orm.Debug(fmt.Sprintf("[backend] entering generateTicker"))
	if tickerMap, e := orm.KlineM1ToTicker(startTS, endTS); e == nil {
		if tickerMap != nil && len(tickerMap) > 0 {
			for k, v := range tickerMap {
				(*latestTicker)[k] = v
			}
			orm.Debug(fmt.Sprintf("returnTickerMap %+v", latestTicker))
		} else {
			orm.Debug(fmt.Sprintf("No product's deal refresh in [%d, %d), latestTicker: %+v", startTS, endTS, *latestTicker))
		}
	} else {
		orm.Error(fmt.Sprintf("generateTicker error %+v, latestTickers %+v, returnTickers: %+v", e, *latestTicker, tickerMap))
		return err
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
