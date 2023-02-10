package sys

import (
	"encoding/json"
	"fmt"
	"github.com/derain/core/logs"
	"github.com/derain/internal/pkg/utils"
	"log"
	"os"
)

// file system table
type TSysLog struct {
	LogItem []any
}

// Initialize the local file system database
func (tf *TSysLog) InitSysLogDB(fileName string) error {
	tsl := TSysLog{
		LogItem: []any{"1", "2"},
	}
	sysDBPath := LoadTSys().DBRootPath + "db/logs/"
	logFilePath := sysDBPath + fileName
	f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE, 0777)
	defer f.Close()
	data, err := json.Marshal(tsl)
	if err != nil {
		fmt.Errorf("error")
	}
	_, err = f.Write(data)
	// create dir
	if !utils.CheckPathExists(logFilePath) {
		os.MkdirAll(logFilePath, 0777)
	}
	return nil
}

func LoadSysLogByType(logType int) *TSysLog {
	dir := LoadTSys().DBRootPath + "db/logs/"
	sysDBPath := dir + "/" + logs.LOG_TYPE_MAPPING[logType]
	fp, err := os.OpenFile(sysDBPath, os.O_RDONLY, 0755)
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 1024*1024)
	n, err := fp.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	var sysLog TSysLog
	err = json.Unmarshal(data[:n], &sysLog)
	return &sysLog
}
