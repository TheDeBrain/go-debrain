package sync

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/derain/core/db/table/sys"
	"github.com/derain/core/protocols"
	"github.com/derain/internal/pkg/rules"
	"io"
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
		protocol := new(protocols.CommProtocol)
		// protocol type handle
		protocolTypeBuf := make([]byte, rules.PROTOCOL_TYPE_BYTE_NUM)
		_, err := conn.Read(protocolTypeBuf)
		if err != nil || err == io.EOF {
			return err
		}
		ptBuf := bytes.NewReader(protocolTypeBuf)
		binary.Read(ptBuf, binary.BigEndian, &protocol.ProtocolType)
		switch protocol.ProtocolType {
		case uint8(rules.FILE_BLOCK_CLIENT_SYNC_PROTOCOL):
			{
				handleClientSyncReq(conn)
				break
			}
		case uint8(rules.FILE_BLOCK_BETWEEN_SERVER_SYNC_PROTOCOL):
			{
				handleBetweenServerSyncReq(conn)
				break
			}
		case uint8(rules.FILE_BLOCK_UPLOAD_SYNC_PROTOCOL):
			{
				handleUploadSyncReq(conn)
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

