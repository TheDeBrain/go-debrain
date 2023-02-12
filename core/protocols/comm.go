package protocols

import (
	"bytes"
	"encoding/binary"
)

// Global communication protocol
type CommProtocol struct {
	ProtocolType uint8 // Declare the protocol type , 1 byte
}

func CPBuf(cp *CommProtocol) (*bytes.Buffer, error) {
	buff := bytes.NewBuffer([]byte{})
	err := binary.Write(buff, binary.BigEndian, cp.ProtocolType)
	if err != nil {
		return nil, err
	}
	return buff, nil
}
