package backend

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/ok-chain/okchain/x/token"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
	"os"
	"runtime/debug"
	"sort"
	"sync"
	"time"
	"github.com/ok-chain/okchain/x/common"
)

// http://gorm.io/docs/query.html
type ORM struct {
	db               *gorm.DB
	logger           *log.Logger
	bufferLock       sync.Locker
	lastK15Timestamp int64
	klineM15sBuffer  map[string][]KlineM15
	lastK1Timestamp  int64
	klineM1sBuffer   map[string][]KlineM1
}

func NewORM(enableLog bool, dbDir string, dbName string, logger *log.Logger) (*ORM, error) {
	orm := ORM{}

	if _, err := os.Stat(dbDir); err != nil {
		if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	orm.Debug(fmt.Sprintf("%s created", dbDir))

	dbPath := dbDir + string(os.PathSeparator) + dbName
	if db, err := gorm.Open("sqlite3", dbPath); err != nil {
		(*logger).Error(fmt.Sprintf("dbPath: %s, error: %+v", dbPath, err))
		panic(err)
	} else {
		orm.logger = logger
		orm.db = db
		orm.bufferLock = new(sync.Mutex)
		orm.db.LogMode(enableLog)
		orm.db.AutoMigrate(&Deal{})
		orm.db.AutoMigrate(&token.FeeDetail{})
		orm.db.AutoMigrate(&Order{})
		orm.db.AutoMigrate(&Transaction{})

		allKlinesMap := GetAllKlineMap()
		for _, v := range *allKlinesMap {
			k, _ := NewKlineFactory(v, nil)
			orm.db.AutoMigrate(k)
		}
	}
	return &orm, nil
}

func (orm *ORM) Close() error {
	return orm.db.Close()
}

func (orm *ORM) AddDeals(deals *[]Deal) (addedCnt int, err error) {
	cnt := 0
	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)

	for _, deal := range *deals {
		ret := tx.Create(&deal)
		if ret.Error != nil {
			return cnt, ret.Error
		} else {
			cnt += 1
		}
	}

	tx.Commit()
	return cnt, nil
}

// Deals
func (orm *ORM) DeleteDealBefore(timestamp int64) (err error) {
	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)

	r := tx.Delete(&Deal{}, " Timestamp < ? ", timestamp)
	if r.Error == nil {
		tx.Commit()
	} else {
		return r.Error
	}
	return nil
}

func (orm *ORM) Debug(msg string) {
	if orm.logger != nil {
		(*orm.logger).Debug(msg)
	} else {
		fmt.Println(msg)
	}
}

//func (orm *ORM) Info(msg string) {
//	if orm.logger != nil {
//		(*orm.logger).Info(msg)
//	} else {
//		fmt.Println(msg)
//	}
//
//}
//
func (orm *ORM) Error(msg string) {
	if orm.logger != nil {
		(*orm.logger).Error(msg)
	} else {
		fmt.Println(msg)
	}

}

func (orm *ORM) deferRollbackTx(trx *gorm.DB, returnErr error)  {
	e := recover()
	if e != nil {
		orm.Error(fmt.Sprintf("Panic : %+v", e))
		debug.PrintStack()
	}
	if e != nil || returnErr != nil {
		trx.Rollback()
	}
}

func (orm *ORM) GetLatestDeals(product string, limit int) (*[]Deal, error) {
	var deals []Deal
	r := orm.db.Where("Product = ?", product).Order("Timestamp desc").Limit(limit).Find(&deals)
	if r.Error != nil {
		return nil, r.Error
	}

	return &deals, r.Error
}

func (orm *ORM) GetDeals(address, product string, offset, limit int) (*[]Deal, int) {
	var deals []Deal
	query := orm.db.Model(Deal{})
	if address != "" {
		query = query.Where("sender = ?", address)
	}
	if product != "" {
		query = query.Where("product = ?", product)
	}
	var total int
	query.Count(&total)
	if offset >= total {
		return &deals, total
	}

	query.Order("timestamp desc").Offset(offset).Limit(limit).Find(&deals)
	return &deals, total
}

func (orm *ORM) GetDealsByTimestampRange(product string, startTS, endTS int64) (*[]Deal, error) {
	var deals []Deal
	r := orm.db.Model(Deal{}).Where(
		"Product = ? and Timestamp >= ? and Timestamp < ?", product, startTS, endTS).Order("Timestamp desc").Find(&deals)
	if r.Error == nil {
		return &deals, nil
	}
	return nil, r.Error
}

func (orm *ORM) getOpenCloseDeals(startTS, endTS int64, product string) (open *Deal, close *Deal) {
	var openDeal, closeDeal Deal
	orm.db.Model(Deal{}).Where("Timestamp >= ? and Timestamp < ? and Product = ?", startTS, endTS, product).Order("Timestamp desc").Limit(1).First(&closeDeal)
	orm.db.Model(Deal{}).Where("Timestamp >= ? and Timestamp < ? and Product = ?", startTS, endTS, product).Order("Timestamp asc").Limit(1).First(&openDeal)

	if startTS <= openDeal.Timestamp && openDeal.Timestamp < endTS {
		return &openDeal, &closeDeal
	}

	return nil, nil
}

func (orm *ORM) getOpenCloseKline(startTS, endTS int64, product string, firstK interface{}, lastK interface{}) error {
	defer common.PrintStackIfPainic()

	orm.db.Where("Timestamp >= ? and Timestamp < ? and Product = ?", startTS, endTS, product).Order("Timestamp desc").Limit(1).First(lastK)
	orm.db.Where("Timestamp >= ? and Timestamp < ? and Product = ?", startTS, endTS, product).Order("Timestamp asc").Limit(1).First(firstK)

	return nil
}

func (orm *ORM) GetDealsMinTimestamp() int64 {

	sql := fmt.Sprintf("select min(Timestamp) as ts from deals")
	ts := int64(-1)

	r := orm.db.Raw(sql).Row()
	if r != nil {
		r.Scan(&ts)
	}

	return ts

}

// KlineM1 GetKlineMaxTimestamp
func (orm *ORM) GetKlineMaxTimestamp(k IKline) int64 {

	sql := fmt.Sprintf("select max(Timestamp) as ts from %s", k.GetTableName())
	ts := int64(-1)

	r := orm.db.Raw(sql).Row()
	if r != nil {
		r.Scan(&ts)
	}
	return ts
}

// KlineM1 GetKlineMaxTimestamp
func (orm *ORM) GetKlineMinTimestamp(k IKline) int64 {

	sql := fmt.Sprintf("select min(Timestamp) as ts from %s", k.GetTableName())
	ts := int64(-1)

	r := orm.db.Raw(sql).Row()
	if r != nil {
		r.Scan(&ts)
	}
	return ts
}

// Rule1. No deals to handle between [startTS, endTS), anchorEndTS <- startTS
func (orm *ORM) Deal2Kline1min(startTS, endTS int64) (anchorEndTS int64, newK int, err error) {

	// 1. Get anchor start time.
	if endTS <= startTS {
		return -1, 0, fmt.Errorf("EndTimestamp %d <= StartTimestamp %d, somewhere goes wrong", endTS, startTS)
	}

	acTS := startTS
	maxTSPersistent := orm.GetKlineMaxTimestamp(&KlineM1{})
	if maxTSPersistent > 0 && maxTSPersistent > startTS {
		acTS = maxTSPersistent
	}

	if startTS == 0 {
		minDealTS := orm.GetDealsMinTimestamp()
		// No Deals to handle if minDealTS == -1, anchorEndTS <-- startTS
		if minDealTS == -1 {
			return startTS, 0, fmt.Errorf("No Deals to handled, return without converting job.")
		} else {
			acTS = minDealTS
		}
	}

	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)

	anchorTime := time.Unix(acTS, 0).UTC()
	anchorStartTime := time.Date(
		anchorTime.Year(), anchorTime.Month(), anchorTime.Day(), anchorTime.Hour(), anchorTime.Minute(), 0, 0, time.UTC)

	// 2. Collect product's kline by deals
	productKlines := map[string][]KlineM1{}
	nextTime := anchorStartTime.Add(time.Minute)
	nextTimeStamp := nextTime.Unix()
	for nextTimeStamp < endTS {

		sql := fmt.Sprintf("select product, sum(Quantity) as quantity, max(Price) as high, min(Price) as low, count(price) as cnt from deals "+
			"where Timestamp >= %d and Timestamp < %d group by product", anchorStartTime.Unix(), nextTime.Unix())

		rows, _ := orm.db.Raw(sql).Rows()

		for rows.Next() {
			var product string
			var quantity, high, low float64
			var cnt int

			rows.Scan(&product, &quantity, &high, &low, &cnt)
			if cnt > 0 {

				openDeal, closeDeal := orm.getOpenCloseDeals(anchorStartTime.Unix(), nextTime.Unix(), product)

				b := BaseKline{
					Product: product, High: high, Low: low, Volume: quantity / 2, Timestamp: anchorStartTime.Unix(),
					Open: openDeal.Price, Close: closeDeal.Price}
				k1min := NewKlineM1(&b)

				klines := productKlines[product]
				if klines == nil {
					klines = []KlineM1{*k1min}
				} else {
					klines = append(klines, *k1min)
				}
				productKlines[product] = klines
			}
		}

		rows.Close()

		anchorStartTime = nextTime
		nextTime = anchorStartTime.Add(time.Minute)
		nextTimeStamp = nextTime.Unix()
	}

	// 3. Batch insert into Kline1Min

	for product, klines := range productKlines {
		for _, kline := range klines {
			// TODO: it should be a replacement here.
			ret := tx.Create(&kline)
			if ret.Error != nil {
				fmt.Printf("Error: %+v, product: %s, kline: %+v", ret.Error, product, kline)
			} else {
				fmt.Printf("%s %+v", timeString(kline.Timestamp), kline)
			}

		}
	}
	tx.Commit()

	anchorEndTS = anchorStartTime.Unix()
	return anchorEndTS, len(productKlines), nil
}

func (orm *ORM) deleteKlinesBefore(unixTS int64, kline interface{}) (err error ){

	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)

	r := tx.Delete(kline, " Timestamp < ? ", unixTS)
	if r.Error == nil {
		tx.Commit()
	} else {
		return r.Error
	}
	return nil
}

func (orm *ORM) DeleteKlineBefore(unixTS int64, kline interface{}) error {
	return orm.deleteKlinesBefore(unixTS, kline)
}

func (orm *ORM) DeleteKlineM1Before(unixTS int64) error {
	return orm.DeleteKlineBefore(unixTS, &KlineM1{})
}

func (orm *ORM) getAllUpdatedProducts(anchorStartTS, anchorEndTS int64) ([]string, error) {
	sql := fmt.Sprintf("select distinct(Product) from deals where Timestamp >= %d and Timestamp < %d",
		anchorStartTS, anchorEndTS)

	rows, err := orm.db.Raw(sql).Rows()

	if err == nil {
		products := []string{}
		for rows.Next() {
			var product string
			rows.Scan(&product)
			products = append(products, product)
		}

		rows.Close()
		return products, nil

	} else {
		return nil, err
	}
}

func (orm *ORM) getLatestKlinesByProduct(product string, limit int, anchorTS int64, klines interface{}) error {

	var r *gorm.DB
	if anchorTS > 0 {
		r = orm.db.Where("Timestamp < ? and Product = ?", anchorTS, product).Order("Timestamp desc").Limit(limit).Find(klines)
	} else {
		r = orm.db.Where("Product = ?", product).Order("Timestamp desc").Limit(limit).Find(klines)
	}

	return r.Error
}

func (orm *ORM) GetKlinesByTimeRange(product string, startTS, endTS int64, klines interface{}) error {

	var r *gorm.DB
	r = orm.db.Where("Timestamp >= ? and Timestamp < ? and Product = ?", startTS, endTS, product).
		Order("Timestamp desc").Find(klines)

	return r.Error
}

func (orm *ORM) GetLatestKlineM1ByProduct(product string, limit int) (*[]KlineM1, error) {
	klines := []KlineM1{}
	if err := orm.getLatestKlinesByProduct(product, limit, -1, &klines); err != nil {
		return nil, err
	} else {
		return &klines, nil
	}
}

// KlineM1 --> KlineM*
func (orm *ORM) MergeKlineM1(startTS, endTS int64, destKline IKline) (anchorEndTS int64, newKCnt int, err error) {

	kM1, _ := NewKlineFactory("kline_m1", nil)
	// 0. destKline should not be KlineM1 & endTS should be greater than startTS
	if destKline.GetFreqInSecond() <= kM1.(IKline).GetFreqInSecond() {
		return startTS, 0, fmt.Errorf("destKline's updating Freq #%d# should be greater than 60", destKline.GetFreqInSecond())
	}
	if endTS <= startTS {
		return -1, 0, fmt.Errorf("EndTimestamp %d <= StartTimestamp %d, somewhere goes wrong", endTS, startTS)
	}

	// 1. Get anchor start time.
	acTS := startTS
	maxTSPersistent := orm.GetKlineMaxTimestamp(destKline)
	if maxTSPersistent > 0 && maxTSPersistent > startTS {
		acTS = maxTSPersistent
	}

	if acTS == 0 {
		minTS := orm.GetKlineMinTimestamp(kM1.(IKline))
		// No Deals to handle if minDealTS == -1, anchorEndTS <-- startTS
		if minTS == -1 {
			return startTS, 0, errors.New("DestKline:" + destKline.GetTableName() + ". No KlineM1 to handled, return without converting job.")
		} else {
			acTS = minTS
		}
	}

	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)

	anchorTime := time.Unix(acTS, 0).UTC()
	var anchorStartTime time.Time
	if maxTSPersistent > 0 {
		anchorStartTime = time.Date(
			anchorTime.Year(), anchorTime.Month(), anchorTime.Day(), anchorTime.Hour(), anchorTime.Minute(), anchorTime.Second(), 0, time.UTC)
	} else {
		anchorStartTime = time.Date(
			anchorTime.Year(), anchorTime.Month(), anchorTime.Day(), anchorTime.Hour(), 0, 0, 0, time.UTC)
	}

	// 2. Get anchor end time.
	anchorEndTime := endTS

	// 3. Collect product's kline by deals
	productKlines := map[string][]interface{}{}
	interval := time.Duration(int(time.Second) * destKline.GetFreqInSecond())
	nextTime := anchorStartTime.Add(interval)
	nextTimeStamp := nextTime.Unix()
	for nextTimeStamp < anchorEndTime {

		sql := fmt.Sprintf("select %d, product, sum(volume) as volume, max(high) as high, min(low) as low, count(*) as cnt from %s "+
			"where Timestamp >= %d and Timestamp < %d group by product", anchorStartTime.Unix(), kM1.(IKline).GetTableName(), anchorStartTime.Unix(), nextTime.Unix())

		rows, _ := orm.db.Raw(sql).Rows()

		for rows.Next() {
			var product string
			var quantity, high, low float64
			var cnt int
			var ts int64

			rows.Scan(&ts, &product, &quantity, &high, &low, &cnt)
			if cnt > 0 {

				openKline, _ := NewKlineFactory(kM1.(IKline).GetTableName(), nil)
				closeKline, _ := NewKlineFactory(kM1.(IKline).GetTableName(), nil)
				orm.getOpenCloseKline(anchorStartTime.Unix(), nextTime.Unix(), product, openKline, closeKline)

				b := BaseKline{
					Product: product, High: high, Low: low, Volume: quantity, Timestamp: anchorStartTime.Unix(),
					Open: openKline.(IKline).GetOpen(), Close: closeKline.(IKline).GetClose()}

				newDestK, _ := NewKlineFactory(destKline.GetTableName(), &b)

				klines := productKlines[product]
				if klines == nil {
					klines = []interface{}{newDestK}
				} else {
					klines = append(klines, newDestK)
				}
				productKlines[product] = klines
			}
		}

		rows.Close()

		anchorStartTime = nextTime
		nextTime = anchorStartTime.Add(interval)
		nextTimeStamp = nextTime.Unix()
	}

	// 4. Batch insert into Kline1Min

	for product, klines := range productKlines {
		for _, kline := range klines {
			// TODO: it should be a replacement here.
			ret := tx.Create(kline)
			if ret.Error != nil {
				fmt.Printf("Error: %+v, product: %s, kline: %+v", ret.Error, product, kline)
			} else {
				fmt.Printf("%s %+v", timeString(kline.(IKline).GetTimestamp()), kline)
			}
		}
	}
	tx.Commit()

	anchorEndTS = anchorStartTime.Unix()
	return anchorEndTS, len(productKlines), nil
}

// Latest 24H KlineM1 to Ticker
func (orm *ORM) KlineM1ToTicker(startTS, endTS int64) (m map[string]*Ticker, err error) {
	orm.bufferLock.Lock()
	defer orm.bufferLock.Unlock()

	// 1. Get updated product by KlineM1 in latest 120 seconds
	km1, _ := NewKlineFactory("kline_m1", nil)
	km15, _ := NewKlineFactory("kline_m15", nil)
	anchorKM1TS := (km1).(IKline).GetAnchorTimeTS(endTS)
	productList, _ := orm.getAllUpdatedProducts(startTS, endTS)
	if productList == nil || len(productList) == 0 {
		return nil, nil
	}

	// 2. Update Buffer.
	// 	2.1 For each product, get latest [anchorKM15TS-95*60, anchorKM15TS) KlineM15 list
	anchorKM15TS := km15.(IKline).GetAnchorTimeTS(endTS)
	bufferKM15 := make(map[string][]KlineM15)
	if anchorKM15TS == orm.lastK15Timestamp {
		bufferKM15 = orm.klineM15sBuffer
	}

	for _, p := range productList {
		existsKM15 := bufferKM15[p]
		if existsKM15 != nil && len(existsKM15) > 0 {
			continue
		}

		klineM15s := []KlineM15{}
		orm.getLatestKlinesByProduct(p, 95, anchorKM15TS, &klineM15s)
		bufferKM15[p] = klineM15s
	}
	orm.lastK15Timestamp = anchorKM15TS
	orm.klineM15sBuffer = bufferKM15

	// 	2.2 For each product, get latest [anchorKM15TS, anchorKM1TS) KlineM1 list
	bufferKM1 := map[string][]KlineM1{}
	if anchorKM1TS == orm.lastK1Timestamp {
		bufferKM1 = orm.klineM1sBuffer
	}

	for _, p := range productList {
		existsKM1 := bufferKM1[p]
		if existsKM1 != nil && len(existsKM1) > 0{
			continue
		}

		klineM1s := []KlineM1{}
		orm.GetKlinesByTimeRange(p, anchorKM15TS, anchorKM1TS, &klineM1s)
		bufferKM1[p] = klineM1s
	}
	orm.lastK1Timestamp = anchorKM1TS
	orm.klineM1sBuffer = bufferKM1

	// 	2.3 For each product, get latest [anchorKM1TS, endTS) Deal list
	bufferDeals := map[string][]Deal{}
	for _, p := range productList {
		deals, _ := orm.GetDealsByTimestampRange(p, anchorKM1TS, endTS)
		if deals != nil && len(*deals) > 0 {
			bufferDeals[p] = *deals
		}
	}

	// 3. For each updated product, generate new ticker by KlineM15 & KlineM1 & Deals in 24 Hours
	orm.Debug(fmt.Sprintf("KlineM15: %+v\n KlineM1: %+v\n Deals: %+v\n", orm.klineM15sBuffer, orm.klineM1sBuffer, bufferDeals))
	tickerMap := map[string]*Ticker{}

	orm.Debug(fmt.Sprintf("KlineM1ToTicker's productList %+v", productList))
	for _, p := range productList {
		klinesM1 := bufferKM1[p]
		klinesM15 := bufferKM15[p]
		iklines := IKlines{}

		for idx, _ := range klinesM1[:] {
			iklines = append(iklines, &klinesM1[idx])
		}
		for idx, _ := range klinesM15[:] {
			iklines = append(iklines, &klinesM15[idx])
		}

		// [X] 3.1 No klinesM1 & klinesM15 found, contine
		// FLT. 20190411. Go ahead even if there's no klines.

		// 3.2 Do iklines sort desc by timestamp.

		allVolume, lowest, highest := 0.0, 0.0, 0.0
		deals := bufferDeals[p]

		if len(iklines) > 0 {
			sort.Sort(iklines)
			fmt.Printf("KlineM1ToTicker IKlines %+v\n", iklines)
			allVolume, lowest, highest = 0.0, iklines[0].GetLow(), iklines[0].GetHigh()
			for idx, k := range iklines {
				fmt.Printf("%d, Handled Kline: %s\n", idx, k.PrettyTimeString())

				allVolume += k.GetVolume()
				if k.GetHigh() > highest {
					highest = k.GetHigh()
				}
				if k.GetLow() < lowest {
					lowest = k.GetLow()
				}
			}
		} else {
			if len(deals) > 0 {
				allVolume, lowest, highest = 0.0, deals[0].Price, deals[0].Price
			}
		}

		fmt.Printf("KlineM1ToTicker Deals %+v\n", deals)
		for _, d := range deals {
			allVolume += (d.Quantity / 2)
			if d.Price > highest {
				highest = d.Price
			}
			if d.Price < lowest {
				lowest = d.Price
			}
		}

		t := Ticker{}
		if len(iklines) > 0 {
			t.Open = iklines[len(iklines)-1].GetOpen()
		} else {
			t.Open = deals[len(deals)-1].Price
		}

		if len(deals) > 0 {
			t.Close = deals[0].Price
		} else {
			t.Close = iklines[0].GetClose()
		}

		t.Volume = allVolume
		t.High = highest
		t.Low = lowest
		t.Symbol = p
		t.CurrencyId = p
		t.Change = (t.Close - t.Open)
		t.ChangePercentage = t.Change / t.Open
		t.Timestamp = endTS
		tickerMap[p] = &t
	}

	for k, v := range tickerMap {
		orm.Debug(fmt.Sprintf("KlineM1ToTicker Ticker[%s] %s", k, v.PrettyString()))
	}
	return tickerMap, nil
}

// FeeDetail
func (orm *ORM) AddFeeDetails(feeDetails *[]token.FeeDetail) (addedCnt int, err error) {
	cnt := 0
	tx := orm.db.Begin()

	for _, feeDetail := range *feeDetails {
		ret := tx.Create(&feeDetail)
		if ret.Error != nil {
			tx.Rollback()
			return cnt, ret.Error
		} else {
			cnt += 1
		}
	}

	tx.Commit()
	return cnt, nil
}

func (orm *ORM) GetFeeDetails(address string, offset, limit int) (*[]token.FeeDetail, int) {
	var feeDetails []token.FeeDetail
	query := orm.db.Model(token.FeeDetail{}).Where("address = ?", address)
	var total int
	query.Count(&total)
	if offset >= total {
		return &feeDetails, total
	}

	query.Order("timestamp desc").Offset(offset).Limit(limit).Find(&feeDetails)
	return &feeDetails, total
}

// Order
func (orm *ORM) AddOrders(orders *[]Order) (addedCnt int, err error) {
	cnt := 0
	tx := orm.db.Begin()

	for _, order := range *orders {
		ret := tx.Create(&order)
		if ret.Error != nil {
			tx.Rollback()
			return cnt, ret.Error
		} else {
			cnt += 1
		}
	}

	tx.Commit()
	return cnt, nil
}

func (orm *ORM) UpdateOrders(orders *[]Order) (int, error) {
	cnt := 0
	tx := orm.db.Begin()

	for _, order := range *orders {
		ret := tx.Save(&order)
		if ret.Error != nil {
			tx.Rollback()
			return cnt, ret.Error
		} else {
			cnt += 1
		}
	}

	tx.Commit()
	return cnt, nil
}

func (orm *ORM) GetOrderList(address, product string, open bool, offset, limit int) (*[]Order, int) {
	var orders []Order
	query := orm.db.Model(Order{}).Where("sender = ?", address)
	if product != "" {
		query = query.Where("product = ?", product)
	}
	if open {
		query = query.Where("status = 0")
	} else {
		query = query.Where("status > 0")
	}
	var total int
	query.Count(&total)
	if offset >= total {
		return &orders, total
	}

	query.Order("timestamp desc").Offset(offset).Limit(limit).Find(&orders)
	return &orders, total
}

// Transaction
func (orm *ORM) AddTransactions(transactions *[]Transaction) (addedCnt int, err error) {
	cnt := 0
	tx := orm.db.Begin()

	for _, transaction := range *transactions {
		ret := tx.Create(&transaction)
		if ret.Error != nil {
			tx.Rollback()
			return cnt, ret.Error
		} else {
			cnt += 1
		}
	}

	tx.Commit()
	return cnt, nil
}

func (orm *ORM) GetTransactionList(address string, txType, startTime, endTime int64, offset, limit int) (*[]Transaction, int) {
	var txs []Transaction
	query := orm.db.Model(Transaction{}).Where("address = ?", address)
	if txType != 0 {
		query = query.Where("type = ?", txType)
	}
	if startTime > 0 {
		query = query.Where("timestamp >= ?", startTime)
	}
	if endTime > 0 {
		query = query.Where("timestamp < ?", endTime)
	}

	var total int
	query.Count(&total)
	if offset >= total {
		return &txs, total
	}

	query.Order("timestamp desc").Offset(offset).Limit(limit).Find(&txs)
	return &txs, total
}