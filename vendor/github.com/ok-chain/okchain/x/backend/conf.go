package backend

import (
	"bytes"
	"encoding/json"
	"github.com/tendermint/tendermint/libs/common"
	"os"
	)

var (
	DefaultMaintainConfile 	= "maintain.conf"
	DefaultNodeHome 		= os.ExpandEnv("$HOME/.okchaind")
	DefaultNodeCofig		= DefaultNodeHome + "/config"
	DefaultNodeDataHome 	= DefaultNodeHome + "/data"
	DefaultTestConfig       = DefaultNodeHome + "/test_config"
	DefaultTestDataHome     = DefaultNodeHome + "/test_data"
)

type MaintainConf struct {
	EnableBackend bool   `json:"enable_backend"`
	HotKeptDays   int    `json:"hot_kept_days"`
	UpdateFreq    int64  `json:"update_freq"`   // unit: second
	BufferSize    int    `json:"buffer_size"`   //
	SyncMode      string `json:"sync_mode"`     // mode: block or minutes
	Sqlite3Path   string `json:"sqlite_3_path"` // path: ~/.okdex/db/sqlite3
	LogSQL        bool   `json:"log_sql"`		//
	CleanUpsKeptDays map[string] int `json:"clean_ups_kept_days"` // 0 <= x <= 60
	CleanUpsTime  string `json:"clean_ups_time"`// e.g.) 00:00:00, CleanUp job will be fired at this time.
	GenesisTime   string `json:"genesis_time"`  // genesis_time: "2019-05-01 00:00:00"
}

func GetDefaultMaintainConfig() *MaintainConf {
	m := MaintainConf{}

	m.EnableBackend = false
	m.HotKeptDays = 3
	m.UpdateFreq = 60
	m.BufferSize = 4096
	m.Sqlite3Path = DefaultNodeDataHome + "/sqlite3"
	m.LogSQL = true
	m.CleanUpsTime = "00:00:00"
	m.CleanUpsKeptDays = map[string]int{}
	m.CleanUpsKeptDays["kline_m1"] = 30
	m.CleanUpsKeptDays["kline_m3"] = 30
	m.CleanUpsKeptDays["kline_m5"] = 30
	m.GenesisTime = "2019-04-01 00:00:00"

	return &m
}

func LoadMaintainConf(confDir string, fileName string) (*MaintainConf, error) {
	fPath := confDir + string(os.PathSeparator) + fileName
	if _, err := os.Stat(fPath); err != nil {
		return nil, err
	}

	bytes := common.MustReadFile(fPath)

	m := MaintainConf{}
	err := json.Unmarshal(bytes, &m)
	return &m, err
}

func DumpMaintainConf(maintainConf *MaintainConf, confDir string, fileName string) error {
	fPath := confDir + string(os.PathSeparator) + fileName

	if _, err := os.Stat(confDir); err != nil {
		os.MkdirAll(confDir, os.ModePerm)
	}

	if bs, err := json.Marshal(maintainConf); err != nil {
		return err
	} else {
		var out bytes.Buffer
		json.Indent(&out, bs, "", "  ")
		common.MustWriteFile(fPath, out.Bytes(), os.ModePerm)
	}

	return nil
}

func SafeLoadMaintainConfig(configDir string) *MaintainConf {
	maintainConf, _ := LoadMaintainConf(configDir, DefaultMaintainConfile)
	if maintainConf == nil {
		maintainConf = GetDefaultMaintainConfig()
		DumpMaintainConf(maintainConf, configDir, DefaultMaintainConfile)
	}
	return maintainConf
}
