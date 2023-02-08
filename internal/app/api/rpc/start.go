package rpc

import (
	sys2 "github.com/derain/core/db/table/sys"
	"github.com/derain/internal/app/api/rpc/services"
	"log"
	"net"
	"net/rpc"
)

func StartRpcService() {
	sys := sys2.LoadTSys()
	ln, err := net.Listen("tcp", ":"+sys.RpcPort)
	log.Println("Start the synchronization server and listen to the port:", sys.RpcPort)
	if err != nil {
		log.Fatal("rpc service error:", err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		defer conn.Close()
		if err != nil {
			log.Fatal("rpc connect error:", err)
		}
		go initRpcService(conn)
	}
}

func initRpcService(conn net.Conn) {
	p := rpc.NewServer()
	p.Register(&services.SysService{Conn: conn})
	p.ServeConn(conn)
}
