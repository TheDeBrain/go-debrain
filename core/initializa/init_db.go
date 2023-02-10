package initializa

import (
	"github.com/derain/core/db/table/node"
	"github.com/derain/core/db/table/sys"
	"github.com/derain/core/logs"
)

func InitDB() error {
	// system
	ts := new(sys.TSys)
	// file system
	fs := new(sys.TFileSys)
	// route table
	rt := new(node.TRouteTable)
	// sys log
	sl := new(sys.TSysLog)
	err := ts.IniSysDB("sys.json", "", "", "", "")
	if err == nil {
		err = fs.InitFileSysDB("file_sys.json")
		err = rt.InitRouteTable("route_table.json")
		// init system log
		err = sl.InitSysLogDB(logs.LOG_TYPE_MAPPING[logs.SYNC_LOG])
	}
	if err != nil {
		return err
	}
	return nil
}

//Initialize the sys log
func InitSysLog() {

}
