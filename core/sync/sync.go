package sync

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/derain/core/db/table/sys"
	"github.com/derain/core/protocols"
	"github.com/derain/core/rules"
	"github.com/derain/core/sync/file"
	"log"
	"net"
)

// start sync service udp
func StartSyncServiceUDP() {
	addr,_:=net.ResolveUDPAddr("udp", ":"+sys.TSysNew().SyncPortUDP)
	listener,err := net.ListenUDP("udp",addr)
	defer listener.Close()
	if err != nil {
		// reconnect
		StartSyncServiceUDP()
	}
	log.Println("Start the rpc server and listen to the udp port:", sys.TSysNew().SyncPortUDP)
	HandleSyncServiceUDP(listener)
}

// handle sync service udp
func HandleSyncServiceUDP(conn *net.UDPConn) {
	for {
		protocol := new(protocols.NetPack)
		// protocol type handle
		protocolTypeBuf := make([]byte, rules.NET_ACTION_TYPE_SIZE)
		_,_,err := conn.ReadFromUDP(protocolTypeBuf)
		if err != nil {
			continue
		}
		ptBuf := bytes.NewReader(protocolTypeBuf)
		binary.Read(ptBuf, binary.BigEndian, &protocol.NetActionType)
		switch protocol.NetActionType {
		case rules.FILE_BLOCK_CLIENT_SYNC_REQ:
			{
				file.HandleClientSyncReqUDP(conn)
				break
			}
		case rules.FILE_BLOCK_CLIENT_SYNC_RECEIVE:
			{
				file.HandleClientSyncReqReceiveUDP(conn)
				break
			}
		case rules.FILE_BLOCK_SERVER_BROADCAST_SYNC:
			{
				file.HandleFileBlockServerBroadCastSyncUDP(conn)
				break
			}
		case rules.FILE_BLOCK_UPLOAD_SYNC_PROTOCOL:
			{
				file.HandleUploadSyncReqUDP(conn)
				break
			}
		case rules.FILE_GETTER_PROTOCOL:
			{
				file.HandleGetFileResponseUDP(conn)
				break
			}
		default:
			{
				continue
			}
		}
	}
}

// start sync service tcp
func StartSyncServiceTCP() error {
	sys := sys.TSysNew()
	ln, err := net.Listen("tcp", ":"+sys.SyncPortTCP)
	defer ln.Close()
	if err != nil {
		return err
	}
	log.Println("Start the rpc server and listen to the tcp port:", sys.SyncPortTCP)
	for {
		conn, _ := ln.Accept()
		if err != nil {
			log.Fatal("sync connect error:", err)
		}
		go HandleSyncServiceTCP(conn)
	}
	return nil
}

// handle sync service
func HandleSyncServiceTCP(conn net.Conn) error {
	for {
		protocol := new(protocols.NetPack)
		// protocol type handle
		protocolTypeBuf := make([]byte, rules.NET_ACTION_TYPE_SIZE)
		_, err := conn.Read(protocolTypeBuf)
		if err != nil {
			return err
		}
		ptBuf := bytes.NewReader(protocolTypeBuf)
		binary.Read(ptBuf, binary.BigEndian, &protocol.NetActionType)
		switch protocol.NetActionType {
		case rules.FILE_BLOCK_CLIENT_SYNC_REQ:
			{
				file.HandleClientSyncReqTCP(conn)
				break
			}
		case rules.FILE_BLOCK_CLIENT_SYNC_RECEIVE:
			{
				file.HandleClientSyncReqReceiveTCP(conn)
				break
			}
		case rules.FILE_BLOCK_SERVER_BROADCAST_SYNC:
			{
				file.HandleFileBlockServerBroadCastSyncTCP(conn)
				break
			}
		case rules.FILE_BLOCK_UPLOAD_SYNC_PROTOCOL:
			{
				file.HandleUploadSyncReqTCP(conn)
				break
			}
		case rules.FILE_GETTER_PROTOCOL:
			{
				file.HandleGetFileResponseTCP(conn)
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
