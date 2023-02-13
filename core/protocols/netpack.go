package protocols

import (
	"bytes"
	"encoding/binary"
	"github.com/derain/core/db/table/node"
	"io"
	"net"
)

// net action protocol
type NetPack struct {
	NetActionType uint8  `json:"net_action_type"` // Declare the net action , 1 byte
	Data          []byte `json:"data"`
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
func NPBuf(np *NetPack) (*bytes.Buffer, error) {
	buff := bytes.NewBuffer([]byte{})
	// read in net action type
	err := binary.Write(buff, binary.BigEndian, np.NetActionType)
	// read in data
	err = binary.Write(buff, binary.BigEndian, np.Data)
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

// net pack full send
func (np *NetPack) NPSendFull(nodeList []node.TNodeInfo) (*ResultCollect, error) {
	var resArr []Result
	for _, n := range nodeList {
		c, err := net.Dial("tcp", n.Addr+":"+n.Port)
		if err != nil {
			// bad node hanlde
		}
		err = NPWriter(c, np)
		if err != nil {
			// bad node hanlde
		}
		// result
		res, err := RESReader(c)
		resArr = append(resArr, *res)
	}
	rc, _ := RCNew(resArr)
	return rc, nil
}

// net pack on send
func (np *NetPack) NPSendOne(n node.TNodeInfo) {
	c, err := net.Dial("tcp", n.Addr+":"+n.Port)
	if err != nil {
		// bad node hanlde
	}
	err = NPWriter(c, np)
	if err != nil {
		// bad node hanlde
	}
}

// ------------------------- struct handle end -------------------------
