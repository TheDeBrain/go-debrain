package initializa

import (
	"github.com/derain/core/db/table/node"
	"github.com/derain/core/db/table/sys"
	"github.com/derain/internal/pkg/vars"
)

func InitDB() error {
	// system
	ts := new(sys.TSys)
	// file system
	fs := new(sys.TFileSys)
	// route table
	rt := new(node.TRouteTable)
	err := ts.IniSysDB("sys.json", "", "", "", "")
	err = fs.InitFileSysDB("file_sys.json")
	err = rt.InitRouteTable("route_table.json")
	vars.TSys = sys.LoadTSys()
	vars.TFSys = sys.LoadFileSys()
	if err != nil {
		return err
	}
	return nil
}

//Initialize the sys log
func InitSysLog() {

}
