package sys

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// system table
type TSys struct {
	SyncPortTCP       string  `json:"sync_port_tcp"` // sync port
	SyncPortUDP       string  `json:"sync_port_udp"`
	RpcPort           string  `json:"rpc_port"`           // rpc port
	WebApiPort        string  `json:"web_api_port"`       // web api port
	HeartbeatInterval int    `json:"heartbeat_interval"` // node heartbeat detection interval, unit: second (s)
	DBRootPath        string `json:"db_root_path"`       // system root path,used to store system data,modification not recommended,Important！！
	Version           string `json:"version"`            // go-debrain version
}

//Initialize the system database
func (ts *TSys) IniSysDB(fileName string,
	syncPortTCP string,
	syncPortUDP string,
	rpcPort string,
	webApiPort string,
	dbRootPath string) error {
	dir, _ := os.Getwd()
	if  len(syncPortTCP)==0 {
		syncPortTCP = "9000"
	}
	if len(syncPortUDP)==0 {
		syncPortUDP = "9001"
	}
	if len(rpcPort)==0 {
		rpcPort = "9002"
	}
	if len(webApiPort)==0 {
		webApiPort = "9003"
	}
	if len(dbRootPath) == 0 {
		dbRootPath = dir + "/debrain-data/"
	}
	dbRootPath = filepath.FromSlash(dbRootPath)
	sysDB := TSys{
		SyncPortTCP:       syncPortTCP,
		SyncPortUDP:       syncPortUDP,
		RpcPort:           rpcPort,
		WebApiPort:        webApiPort,
		HeartbeatInterval: 1,
		DBRootPath:        dbRootPath,
		Version:           "v0.0.1",
	}
	sysDBPath := dir + "/" + fileName
	f, err := os.OpenFile(sysDBPath, os.O_RDWR|os.O_CREATE, 0777)
	defer f.Close()
	data, err := json.Marshal(sysDB)
	if err != nil {
		fmt.Errorf("error")
	}
	_, err = f.Write(data)
	return nil
}

func TSysNew() *TSys {
	dir, _ := os.Getwd()
	sysDBPath := dir + "/sys.json"
	fp, err := os.OpenFile(sysDBPath, os.O_RDONLY, 0777)
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 1024*1024)
	n, err := fp.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	var sys TSys
	err = json.Unmarshal(data[:n], &sys)
	return &sys
}
