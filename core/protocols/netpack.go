package protocols

import (
	"bytes"
	"encoding/binary"
	"io"
)

// net action protocol
type NetPack struct {
	NetActionType uint8 // Declare the net action , 1 byte
	Data          any
}

// ------------------------- struct handle start -------------------------

// create file block pointer
func NPNew(netActionType uint8, data []byte) *NetPack {
	na := NetPack{
		netActionType,
		data,
	}
	return &na
}

// create file block buffer
func NPBuf(na *NetPack) (*bytes.Buffer, error) {
	buff := bytes.NewBuffer([]byte{})
	// read in net action type
	err := binary.Write(buff, binary.BigEndian, na.NetActionType)
	// read in data
	err = binary.Write(buff, binary.BigEndian, na.Data)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

// file block protocol writer
func NPWriter(w io.Writer, na *NetPack) error {
	naArr, err := NPBuf(na)
	if err != nil {
		return err
	}
	w.Write(naArr.Bytes())
	return nil
}

// ------------------------- struct handle end -------------------------
