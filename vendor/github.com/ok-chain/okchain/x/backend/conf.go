package backend

import (
	"bytes"
	"encoding/json"
	"github.com/tendermint/tendermint/libs/common"
	"os"
)

type MaintainConf struct {
	EnableBackend bool   `json:"enable_backend"`
	HotKeptDays   int    `json:"hot_kept_days"`
	UpdateFreq    int64  `json:"update_freq"`   // unit: second
	BufferSize    int    `json:"buffer_size"`   //
	SyncMode      string `json:"sync_mode"`     // mode: block or minutes
	Sqlite3Path   string `json:"sqlite_3_path"` // path: ~/.okdex/db/sqlite3
	LogSQL        bool   `json:"log_sql"`
	CleanUpHour   int    `json:"clean_up_hour"` // 0 <= h <= 23
	GenesisTime   string `json:"genesis_time"`  // genesis_time: "2019-05-01 00:00:00"
}

func GetDefaultMaintainConfig() *MaintainConf {

	m := MaintainConf{}

	m.EnableBackend = false
	m.HotKeptDays = 3
	m.UpdateFreq = 60
	m.BufferSize = 4096
	m.Sqlite3Path = "/tmp/sqlite3"
	m.LogSQL = true
	m.CleanUpHour = 1
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
	if err := json.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func DumpMaintainConf(maintainConf *MaintainConf, confDir string, fileName string) error {
	fPath := confDir + string(os.PathSeparator) + fileName

	if bs, err := json.Marshal(maintainConf); err != nil {
		return err
	} else {
		var out bytes.Buffer
		err = json.Indent(&out, bs, "", "  ")
		if err == nil {
			common.MustWriteFile(fPath, out.Bytes(), os.ModePerm)
		} else {
			return err
		}
	}

	return nil
}
