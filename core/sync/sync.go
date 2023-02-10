package sync

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/derain/core/protocols"
	"github.com/derain/internal/pkg/rules"
	"io"
	"net"
)

// handle sync service
func handleSyncService(conn net.Conn) error {
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

