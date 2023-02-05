package io

import (
	"encoding/json"
	"github.com/derain/core/db/table/sys"
	"log"
	"os"
)

// get system db
func GetSysDB() sys.TSys {
	dir, _ := os.Getwd()
	sysDBPath := dir + "/sys.json"
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
	var sys sys.TSys
	err = json.Unmarshal(data[:n], &sys)
	return sys
}

// get file system db
func GetFileSysDB() sys.TFileSys {
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
	var sys sys.TFileSys
	err = json.Unmarshal(data[:n], &sys)
	return sys
}
