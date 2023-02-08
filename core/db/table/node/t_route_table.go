package node

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type TRouteTable struct {
	NodeList []TNodeInfo `json:"node_list"`
}

//Initialize the route table
func (rt *TRouteTable) InitRouteTable(fileName string) error {
	dir, _ := os.Getwd()
	routeTable := TRouteTable{}
	// test
	routeTable.NodeList = append(routeTable.NodeList, TNodeInfo{"127.0.0.1", "9090"})
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

func (rt *TRouteTable) Load() *TRouteTable {
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
