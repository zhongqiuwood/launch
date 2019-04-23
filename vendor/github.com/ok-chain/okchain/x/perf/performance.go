package perf

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"sync"
	"time"
)

// http://gitlab.okcoin-inc.com/dex/okchain/issues/32

const (
	orderModule        = "order"
	tokenModule        = "token"
	stakingModule      = "staking"
	govModule          = "gov"
	distributionModule = "distribution"

	appFormat = "OKChain block height, %d, BeginBlock elapsed, %d, DeliverTx elapsed, %d, TxNum, %d, EndBlock elapsed, %d, Commit elapsed, %d,"
	moduleFormat = "OKChain block height, %d, module, %s, BeginBlock elapsed, %d, DeliverTx elapsed, %d, TxNum, %d, EndBlock elapsed, %d,"
	handlerFormat = "OKChain block height, %d, module, %s, handler, %s, elapsed, %d, invoked, %d,"
)
var perf *performance
var once sync.Once

func GetPerf() Perf {
	once.Do(func() {
		perf = newPerf()
	})
	return perf
}

type Perf interface {

	OnAppBeginBlockEnter(height int64) uint64
	OnAppBeginBlockExit(height int64, seq uint64)

	OnAppEndBlockEnter(height int64) uint64
	OnAppEndBlockExit(height int64, seq uint64)

	OnCommitEnter(height int64) uint64
	OnCommitExit(height int64, seq uint64, logger log.Logger)

	OnBeginBlockEnter(ctx sdk.Context, moduleName string) uint64
	OnBeginBlockExit(ctx sdk.Context, moduleName string, seq uint64)

	OnDeliverTxEnter(ctx sdk.Context, moduleName, handlerName string) uint64
	OnDeliverTxExit(ctx sdk.Context, moduleName, handlerName string, seq uint64)

	OnEndBlockEnter(ctx sdk.Context, moduleName string) uint64
	OnEndBlockExit(ctx sdk.Context, moduleName string, seq uint64)

}

type hanlderInfo struct {
	invoke uint64
	elapse int64
}

type info struct {
	blockheight int64
	beginBlockElapse int64
	endBlockElapse int64
	deliverTxElapse int64
	txNum uint64
}

type moduleInfo struct {
	info
	data handlerInfoMapType
}

type appInfo struct {
	info
	commitElapse int64
	lastTimestamp int64
	seqNum uint64
}

type handlerInfoMapType map[string]*hanlderInfo

func newHanlderMetrics() *moduleInfo {
	m := &moduleInfo{
		data: make(handlerInfoMapType),
	}
	return m
}

type performance struct {

	lastTimestamp int64
	seqNum uint64

	app *appInfo
	moduleInfoMap map[string]*moduleInfo
}

func newPerf() *performance {
	p := &performance{
		moduleInfoMap: make(map[string]*moduleInfo),
	}

	p.app = &appInfo{}
	p.moduleInfoMap[orderModule] = newHanlderMetrics()
	p.moduleInfoMap[tokenModule] = newHanlderMetrics()
	p.moduleInfoMap[govModule] = newHanlderMetrics()
	p.moduleInfoMap[distributionModule] = newHanlderMetrics()
	p.moduleInfoMap[stakingModule] = newHanlderMetrics()
	return p
}
////////////////////////////////////////////////////////////////////////////////////

func (p *performance) OnAppBeginBlockEnter(height int64) uint64 {
	p.app.blockheight = height
	p.app.seqNum++
	p.app.lastTimestamp = time.Now().UnixNano()

	return p.app.seqNum
}

func (p *performance) OnAppBeginBlockExit(height int64, seq uint64) {
	p.sanityCheckApp(height, seq)
	p.app.beginBlockElapse = time.Now().UnixNano() - p.app.lastTimestamp
}

////////////////////////////////////////////////////////////////////////////////////


func (p *performance) OnAppEndBlockEnter(height int64) uint64 {
	p.sanityCheckApp(height, p.app.seqNum)

	p.app.seqNum++
	p.app.lastTimestamp = time.Now().UnixNano()

	return p.app.seqNum
}

func (p *performance) OnAppEndBlockExit(height int64, seq uint64) {
	p.sanityCheckApp(height, seq)
	p.app.endBlockElapse = time.Now().UnixNano() - p.app.lastTimestamp
}

////////////////////////////////////////////////////////////////////////////////////

func (p *performance) OnBeginBlockEnter(ctx sdk.Context, moduleName string) uint64 {
	p.lastTimestamp = time.Now().UnixNano()
	p.seqNum++

	m := p.getModule(moduleName)
	m.blockheight = ctx.BlockHeight()

	return p.seqNum
}

func (p *performance) OnBeginBlockExit(ctx sdk.Context, moduleName string, seq uint64) {
	p.sanityCheck(ctx, seq)
	m := p.getModule(moduleName)
	m.beginBlockElapse = time.Now().UnixNano() - p.lastTimestamp
}

////////////////////////////////////////////////////////////////////////////////////
func (p *performance) OnEndBlockEnter(ctx sdk.Context, moduleName string) uint64 {
	p.lastTimestamp = time.Now().UnixNano()
	p.seqNum++

	return p.seqNum
}


func (p *performance) OnEndBlockExit(ctx sdk.Context, moduleName string, seq uint64) {
	p.sanityCheck(ctx, seq)
	m := p.getModule(moduleName)

	m.endBlockElapse = time.Now().UnixNano() - p.lastTimestamp
}

////////////////////////////////////////////////////////////////////////////////////

func (p *performance) OnDeliverTxEnter(ctx sdk.Context, moduleName, handlerName string) uint64 {

	m := p.getModule(moduleName)
	m.blockheight = ctx.BlockHeight()

	info, ok := m.data[handlerName]
	if !ok {
		info = &hanlderInfo{}
		m.data[handlerName] = info
	}

	p.lastTimestamp = time.Now().UnixNano()
	p.seqNum++
	return p.seqNum
}

func (p *performance) OnDeliverTxExit(ctx sdk.Context, moduleName, handlerName string, seq uint64) {
	p.sanityCheck(ctx, seq)

	m := p.getModule(moduleName)

	info, ok := m.data[handlerName]
	if !ok {
		panic("")
	}
	info.invoke++
	info.elapse = time.Now().UnixNano() - p.lastTimestamp

	m.txNum++
	m.deliverTxElapse += info.elapse

	p.app.txNum++
	p.app.deliverTxElapse += info.elapse
}


////////////////////////////////////////////////////////////////////////////////////

func (p *performance) OnCommitEnter(height int64) uint64 {
	p.sanityCheckApp(height, p.app.seqNum)

	p.app.lastTimestamp = time.Now().UnixNano()
	p.app.seqNum++
	return p.app.seqNum
}

func (p *performance) OnCommitExit(height int64, seq uint64, logger log.Logger) {
	p.sanityCheckApp(height, seq)

	if p.app.txNum == 0 {
		return
	}

	p.app.commitElapse = time.Now().UnixNano() - p.app.lastTimestamp

	logger.Info(fmt.Sprintf(appFormat,
		p.app.blockheight,
		p.app.beginBlockElapse/1e9,
		p.app.deliverTxElapse/1e9,
		p.app.txNum,
		p.app.endBlockElapse/1e9,
		p.app.commitElapse/1e9,
	))

	for moduleName, m := range p.moduleInfoMap {

		if m.blockheight == 0 {
			continue
		}

		logger.Info(fmt.Sprintf(moduleFormat,
			m.blockheight,
			moduleName,
			m.beginBlockElapse/1e9,
			m.deliverTxElapse/1e9,
			m.txNum,
			m.endBlockElapse/1e9,
		))

		for hanlderName, info := range m.data {
			logger.Info(fmt.Sprintf(handlerFormat,
				m.blockheight,
				moduleName,
				hanlderName,
				info.elapse/1e9,
				info.invoke,
			))
		}
	}
	p.app = &appInfo{seqNum : p.app.seqNum,}
	p.moduleInfoMap[orderModule] = newHanlderMetrics()
	p.moduleInfoMap[tokenModule] = newHanlderMetrics()
	p.moduleInfoMap[govModule] = newHanlderMetrics()
	p.moduleInfoMap[distributionModule] = newHanlderMetrics()
	p.moduleInfoMap[stakingModule] = newHanlderMetrics()
}
////////////////////////////////////////////////////////////////////////////////////


func (p *performance) sanityCheck(ctx sdk.Context, seq uint64) {
	if seq != p.seqNum {
		panic("")
	}

	if ctx.BlockHeight() != p.app.blockheight {
		panic("")
	}
}

func (p *performance) sanityCheckApp(height int64, seq uint64) {
	if seq != p.app.seqNum {
		panic("")
	}

	if height != p.app.blockheight {
		panic("")
	}
}

func (p *performance) getModule(moduleName string) *moduleInfo{

	v, ok := p.moduleInfoMap[moduleName]
	if !ok {
		panic("")
	}

	return v
}

