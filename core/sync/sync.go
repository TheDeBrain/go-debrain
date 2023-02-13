package sync

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/derain/core/db/table/sys"
	"github.com/derain/core/protocols"
	"github.com/derain/core/rules"
	"log"
	"net"
)

// start sync service
func StartSyncService() error {
	sys := sys.LoadTSys()
	ln, err := net.Listen("tcp", ":"+sys.SyncPort)
	defer ln.Close()
	if err != nil {
		return err
	}
	log.Println("Start the rpc server and listen to the port:", sys.SyncPort)
	for {
		conn, _ := ln.Accept()
		if err != nil {
			log.Fatal("sync connect error:", err)
		}
		go HandleSyncService(conn)
	}
	return nil
}

// handle sync service
func HandleSyncService(conn net.Conn) error {
	for {
		protocol := new(protocols.NetPack)
		// protocol type handle
		protocolTypeBuf := make([]byte, rules.NET_ACTION_TYPE_SIZE)
		_, err := conn.Read(protocolTypeBuf)
		if err != nil{
			return err
		}
		ptBuf := bytes.NewReader(protocolTypeBuf)
		binary.Read(ptBuf, binary.BigEndian, &protocol.NetActionType)
		switch protocol.NetActionType {
		case rules.FILE_BLOCK_CLIENT_SYNC_REQ:
			{
				handleClientSyncReq(conn)
				break
			}
		case rules.FILE_BLOCK_CLIENT_SYNC_RECEIVE:
			{
				handleClientSyncReqReceive(conn)
				break
			}
		case rules.FILE_BLOCK_SERVER_BROADCAST_SYNC:
			{
				handleFileBlockServerBroadCastSync(conn)
				break
			}
		case rules.FILE_BLOCK_UPLOAD_SYNC_PROTOCOL:
			{
				handleUploadSyncReq(conn)
				break
			}
		case rules.GET_FILE_PROTOCOL:
			{
				// Process file acquisition requests and output file blocks to the client
				handleGetFileResponse(conn)
				break
			}
		default:
			{
				return errors.New("Illegal protocol")
			}
		}
	}
	return nil
}
