package sys

import (
	"encoding/json"
	"fmt"
	"github.com/derain/internal/pkg/utils"
	"log"
	"os"
)

// file system table
type TFileSys struct {
	SysLogPath      string `json:"sys_log_path"`
	FileStoragePath string `json:"file_storage_path"`
}

// Initialize the local file system database
func (tf *TFileSys) InitFileSysDB(fileName string) error {
	dir, _ := os.Getwd()
	tfs := TFileSys{
		SysLogPath: TSysNew().DBRootPath + "logs/",
		FileStoragePath: TSysNew().DBRootPath + "files/",
	}
	sysDBPath := dir + "/" + fileName
	f, err := os.OpenFile(sysDBPath, os.O_RDWR|os.O_CREATE, 0777)
	defer f.Close()
	data, err := json.Marshal(tfs)
	if err != nil {
		fmt.Errorf("error")
	}
	_, err = f.Write(data)
	// create dir
	if !utils.CheckPathExists(tfs.FileStoragePath) {
		os.MkdirAll(tfs.FileStoragePath, 0777)
	}
	return nil
}

func LoadFileSys() *TFileSys {
	dir, _ := os.Getwd()
	sysDBPath := dir + "/file_sys.json"
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
	var sys TFileSys
	err = json.Unmarshal(data[:n], &sys)
	return &sys
}
