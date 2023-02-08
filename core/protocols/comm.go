package protocols

import (
	"bytes"
	"encoding/binary"
)

// Global communication protocol
type CommProtocol struct {
	ProtocolType uint8 // Declare the protocol type , 1 byte
	Data         any   // Protocol Data Body
}

func CPBuf(buff *bytes.Buffer, cp CommProtocol) {
	binary.Write(buff, binary.BigEndian, cp.ProtocolType)
}