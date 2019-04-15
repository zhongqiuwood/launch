package backend

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

type IKline interface {
	GetFreqInSecond() int
	GetAnchorTimeTS(ts int64) int64
	GetTableName() string
	GetProduct() string
	GetTimestamp() int64
	GetOpen() float64
	GetClose() float64
	GetHigh() float64
	GetLow() float64
	GetVolume() float64
	PrettyTimeString() string
}

type IKlines []IKline

func (klines IKlines) Len() int {
	return len(klines)
}

func (c IKlines) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (klines IKlines) Less(i, j int) bool {
	return klines[i].GetTimestamp() > klines[j].GetTimestamp()
}

type BaseKline struct {
	Product   string  `gorm:"PRIMARY_KEY;type:varchar(20)" json:"product"`
	Timestamp int64   `gorm:"PRIMARY_KEY;type:int64" json:"timestamp"`
	Open      float64 `gorm:"type:DOUBLE" json:"open"`
	Close     float64 `gorm:"type:DOUBLE" json:"close"`
	High      float64 `gorm:"type:DOUBLE" json:"high"`
	Low       float64 `gorm:"type:DOUBLE" json:"low"`
	Volume    float64 `gorm:"type:DOUBLE" json:"volume"`
	impl      IKline
}

func (b *BaseKline) GetFreqInSecond() int {
	if b.impl != nil {
		return b.impl.GetFreqInSecond()
	} else {
		return -1
	}
}

func (b *BaseKline) GetTableName() string {
	if b.impl != nil {
		return b.impl.GetTableName()
	} else {
		return "base_kline"
	}
}

func (b *BaseKline) GetAnchorTimeTS(ts int64) int64 {
	m := (ts / int64(b.GetFreqInSecond())) * int64(b.GetFreqInSecond())
	return m
}

func (b *BaseKline) GetProduct() string {
	return b.Product
}

func (b *BaseKline) GetTimestamp() int64 {
	return b.Timestamp
}

func (b *BaseKline) GetOpen() float64 {
	return b.Open
}

func (b *BaseKline) GetClose() float64 {
	return b.Close
}

func (b *BaseKline) GetHigh() float64 {
	return b.High
}

func (b *BaseKline) GetLow() float64 {
	return b.Low
}

func (b *BaseKline) GetVolume() float64 {
	return b.Volume
}

func timeString(ts int64) string {
	return time.Unix(ts, 0).Local().Format("2006-01-02 15:04:05")
}

func (b *BaseKline) PrettyTimeString() string {
	return fmt.Sprintf("Product: %s, Freq: %d, Time: %s, OCHLV(%.4f, %.4f, %.4f, %.4f, %.4f)",
		b.Product, b.GetFreqInSecond(), timeString(b.Timestamp), b.Open, b.Close, b.High, b.Low, b.Volume)
}

type KlineM240 BaseKline

type KlineM1 struct {
	*BaseKline
}

func NewKlineM1(b *BaseKline) *KlineM1 {
	k := KlineM1{b}
	k.impl = &k
	return &k
}

func (k *KlineM1) GetFreqInSecond() int {
	return 60
}

func (k *KlineM1) GetTableName() string {
	return "kline_m1"
}

type KlineM3 struct {
	*BaseKline
}

func NewKlineM3(b *BaseKline) *KlineM3 {
	k := KlineM3{b}
	k.impl = &k
	return &k
}
func (k *KlineM3) GetTableName() string {
	return "kline_m3"
}

func (k *KlineM3) GetFreqInSecond() int {
	return 60 * 3
}

type KlineM5 struct {
	*BaseKline
}

func NewKlineM5(b *BaseKline) *KlineM5 {
	k := KlineM5{b}
	k.impl = &k
	return &k
}
func (k *KlineM5) GetTableName() string {
	return "kline_m5"
}

func (k *KlineM5) GetFreqInSecond() int {
	return 60 * 5
}

type KlineM15 struct {
	*BaseKline
}

func NewKlineM15(b *BaseKline) *KlineM15 {
	k := KlineM15{b}
	k.impl = &k
	return &k
}
func (k *KlineM15) GetTableName() string {
	return "kline_m15"
}

func (k *KlineM15) GetFreqInSecond() int {
	return 60 * 15
}

type KlineM30 struct {
	*BaseKline
}

func NewKlineM30(b *BaseKline) *KlineM30 {
	k := KlineM30{b}
	k.impl = &k
	return &k
}
func (k *KlineM30) GetTableName() string {
	return "kline_m30"
}

func (k *KlineM30) GetFreqInSecond() int {
	return 60 * 30
}

type KlineM60 struct {
	*BaseKline
}

func NewKlineM60(b *BaseKline) *KlineM60 {
	k := KlineM60{b}
	k.impl = &k
	return &k
}
func (k *KlineM60) GetTableName() string {
	return "kline_m60"
}
func (k *KlineM60) GetFreqInSecond() int {
	return 60 * 60
}

type KlineM120 struct {
	*BaseKline
}

func NewKlineM120(b *BaseKline) *KlineM120 {
	k := KlineM120{b}
	k.impl = &k
	return &k
}
func (k *KlineM120) GetTableName() string {
	return "kline_m120"
}
func (k *KlineM120) GetFreqInSecond() int {
	return 60 * 120
}

type KlineM1440 struct {
	*BaseKline
}

func NewKlineM1440(b *BaseKline) *KlineM1440 {
	k := KlineM1440{b}
	k.impl = &k
	return &k
}
func (k *KlineM1440) GetTableName() string {
	return "kline_m1440"
}

func (k *KlineM1440) GetFreqInSecond() int {
	return 60 * 24
}

func NewKlineFactory(name string, baseK *BaseKline) (r interface{}, err error) {
	b := baseK
	if b == nil {
		b = &BaseKline{}
	}

	if name == "kline_m1" {
		return NewKlineM1(b), nil
	}

	if name == "kline_m3" {
		return NewKlineM3(b), nil
	}

	if name == "kline_m5" {
		return NewKlineM5(b), nil
	}

	if name == "kline_m15" {
		return NewKlineM15(b), nil
	}

	if name == "kline_m30" {
		return NewKlineM30(b), nil
	}

	if name == "kline_m60" {
		return NewKlineM60(b), nil
	}

	if name == "kline_m120" {
		return NewKlineM120(b), nil
	}

	if name == "kline_m1440" {
		return NewKlineM1440(b), nil
	}

	return nil, errors.New("No kline constructor function found.")
}

func GetAllKlineMap() *map[int]string {

	m := map[int]string{
		180:   "kline_m3",
		300:   "kline_m5",
		900:   "kline_m15",
		1800:  "kline_m30",
		3600:  "kline_m60",
		7200:  "kline_m120",
		86400: "kline_m1440",
	}
	return &m

}

func GetKlineTableNameByFreq(freq int) string {
	m := GetAllKlineMap()
	name := (*m)[freq]
	return name

}

func NewKlinesFactory(name string) (r interface{}, err error) {

	if name == "kline_m1" {
		return &[]KlineM1{}, nil
	}

	if name == "kline_m3" {
		return &[]KlineM3{}, nil
	}

	if name == "kline_m5" {
		return &[]KlineM5{}, nil
	}

	if name == "kline_m15" {
		return &[]KlineM15{}, nil
	}

	if name == "kline_m30" {
		return &[]KlineM30{}, nil
	}

	if name == "kline_m60" {
		return &[]KlineM60{}, nil
	}

	if name == "kline_m120" {
		return &[]KlineM120{}, nil
	}

	if name == "kline_m1440" {
		return &KlineM1440{}, nil
	}

	return nil, errors.New("No klines constructor function found.")
}
