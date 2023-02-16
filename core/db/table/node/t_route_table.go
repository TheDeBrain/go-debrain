package node

import (
	"encoding/json"
	"fmt"
	"github.com/derain/core/db/table/sys"
	"log"
	"math/rand"
	"os"
	"time"
)

type TRouteTable struct {
	NodeListTCP []TNodeInfo `json:"node_list_tcp"`
	NodeListUDP []TNodeInfo `json:"node_list_udp"`
}

// get route table
func TRTNew() *TRouteTable {
	dir, _ := os.Getwd()
	tRouteTableDBPath := dir + "/route_table.json"
	fp, err := os.OpenFile(tRouteTableDBPath, os.O_RDONLY, 0755)
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 1024*1024)
	n, err := fp.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	var tRouteTable TRouteTable
	err = json.Unmarshal(data[:n], &tRouteTable)
	return &tRouteTable
}

// node random getter
func RandomNodeGetter(random int, netType string) []TNodeInfo {
	var nis []TNodeInfo
	switch netType {
	case "tcp":
		{
			nis = TRTNew().NodeListTCP
		}
	case "udp":
		{
			nis = TRTNew().NodeListUDP
		}
	default:
		{
			nis = TRTNew().NodeListTCP
		}
	}
	if len(nis) <= random || random <= 0 {
		return nis
	} else {
		// random...
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(nis), func(i, j int) { nis[i], nis[j] = nis[j], nis[i] })
		nis = nis[:random]
	}
	return nis
}

//Initialize the route table
func (rt *TRouteTable) InitRouteTable(fileName string) error {
	dir, _ := os.Getwd()
	routeTable := TRouteTable{}
	// loacal address
	routeTable.NodeListTCP = append(routeTable.NodeListTCP, TNodeInfo{"127.0.0.1", sys.TSysNew().SyncPortTCP})
	routeTable.NodeListUDP = append(routeTable.NodeListUDP, TNodeInfo{"127.0.0.1", sys.TSysNew().SyncPortUDP})
	routeTableDBPath := dir + "/" + fileName
	f, err := os.OpenFile(routeTableDBPath, os.O_RDWR|os.O_CREATE, 0777)
	defer f.Close()
	data, err := json.Marshal(routeTable)
	if err != nil {
		fmt.Errorf("error")
	}
	_, err = f.Write(data)
	return nil
}
