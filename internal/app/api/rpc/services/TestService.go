package services

import (
	"encoding/json"
	"github.com/derain/core/db/table/sys"
	"net"
)

type SysService struct {
	Conn net.Conn
}

// get system info
func (p *SysService) GetSysInfo(request string, reply *string) error {
	s, _ := json.MarshalIndent(sys.LoadTSys(), "", " ")
	*reply = string(s)
	return nil
}

// get file system info
func (p *SysService) GetFileSysInfo(request string, reply *string) error {
	s, _ := json.MarshalIndent(sys.LoadFileSys(), "", " ")
	*reply = string(s)
	return nil
}
