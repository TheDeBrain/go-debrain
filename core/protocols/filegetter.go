package protocols

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/derain/core/db/table/node"
	"github.com/derain/core/rules"
	"io"
	"net"
	"unsafe"
)

// ------------------------- struct handle start -------------------------

type FileGetter struct {
	ProtocolType    uint8 `json:"protocol_type"`// Declare the protocol type , 1 byte
	FileOwnerSize   uint32 `json:"file_owner_size"`
	FileNameSize    uint64 `json:"file_name_size"`
	FileEndFlagSize uint32 `json:"file_end_flag_size"`// File block end flag data size , 4 bytz
	FileName        []byte `json:"file_name"`// file name
	FileOwner       []byte `json:"file_owner"`// file owner data
	EndFlag         []byte `json:"end_flag"`// file block end flag
}

func FGNew(ows uint32, fns uint64, fn []byte, fo []byte) *FileGetter {
	fg := FileGetter{
		rules.FILE_GETTER_PROTOCOL,
		ows,
		fns,
		uint32(len([]byte(rules.FILE_BLCOK_END_FLAG))),
		fn,
		fo,
		[]byte(rules.FILE_BLCOK_END_FLAG),
	}
	return &fg
}

func FGBuf(fg *FileGetter) (*bytes.Buffer, error) {
	buff := bytes.NewBuffer([]byte{})
	// read in protocol type
	binary.Write(buff, binary.BigEndian, fg.ProtocolType)
	// read in file index
	binary.Write(buff, binary.BigEndian, fg.FileOwnerSize)
	binary.Write(buff, binary.BigEndian, fg.FileNameSize)
	binary.Write(buff, binary.BigEndian, fg.FileEndFlagSize)
	binary.Write(buff, binary.BigEndian, fg.FileName)
	binary.Write(buff, binary.BigEndian, fg.FileOwner)
	binary.Write(buff, binary.BigEndian, fg.EndFlag)
	return buff, nil
}

// file getter net un package
func FGNetUnPack(conn net.Conn) (*FileGetter, error) {
	fg := new(FileGetter)
	fg, err := FGReader(conn)
	if err != nil {
		return nil, err
	}
	return fg, err
}

// file getter protocol reader
func FGReader(r io.Reader) (*FileGetter, error) {
	fg := new(FileGetter)
	// ---------------------------- protocol head ----------------------------
	// protocol type
	protocolTypeBuf := make([]byte, int(unsafe.Sizeof(fg.ProtocolType)))
	_, err := r.Read(protocolTypeBuf)
	ptBuf := bytes.NewReader(protocolTypeBuf)
	binary.Read(ptBuf, binary.BigEndian, &fg.ProtocolType)
	if err != nil {
		return fg, err
	}
	// file owner size
	fileOwnerSizeBuf := make([]byte, int(unsafe.Sizeof(fg.FileOwnerSize)))
	_, err = r.Read(fileOwnerSizeBuf)
	fos := bytes.NewReader(fileOwnerSizeBuf)
	binary.Read(fos, binary.BigEndian, &fg.FileOwnerSize)
	if err != nil {
		return fg, err
	}
	// file name size
	fileNameSizeBuf := make([]byte, int(unsafe.Sizeof(fg.FileNameSize)))
	_, err = r.Read(fileNameSizeBuf)
	fnsBuf := bytes.NewReader(fileNameSizeBuf)
	binary.Read(fnsBuf, binary.BigEndian, &fg.FileNameSize)
	if err != nil {
		return fg, err
	}
	//  end size
	endBuf := make([]byte, int(unsafe.Sizeof(fg.FileEndFlagSize)))
	_, err = r.Read(endBuf)
	ebBuf := bytes.NewReader(endBuf)
	binary.Read(ebBuf, binary.BigEndian, &fg.FileEndFlagSize)
	if err != nil {
		return fg, err
	}
	// ---------------------------- protocol body ----------------------------
	// file owner data
	fileOwner := make([]byte, fg.FileOwnerSize)
	_, err = r.Read(fileOwner)
	fg.FileOwner = fileOwner
	if err != nil {
		return fg, err
	}
	// file name data
	fileNameBuf := make([]byte, fg.FileNameSize)
	_, err = r.Read(fileNameBuf)
	fg.FileName = fileNameBuf
	if err != nil {
		return fg, err
	}
	// end data
	endB := make([]byte, fg.FileEndFlagSize)
	_, err = r.Read(endB)
	fg.EndFlag = endB
	if err != nil {
		return fg, err
	}
	if string(endB) == rules.FILE_BLCOK_END_FLAG {
		return fg, nil
	}
	return fg, errors.New("illegal agreement")
}

func FGWriter(w io.Writer, fg *FileGetter) error {
	fgBuff, err := FGBuf(fg)
	if err != nil {
		return err
	}
	w.Write(fgBuff.Bytes())
	return nil
}

// ------------------------- struct handle end -------------------------

// write file getter to route table by fileblock protocol
func WFGToRT(
	fg *FileGetter,
	c chan FileBlockSyncResult) error {
	// file block result
	fBr := new(FileBlockSyncResult)
	// bad node
	var badNodeList []any
	// route table
	rt := node.TRTNew()
	// file block sync
	for _, n := range rt.NodeListTCP {
		c, er := net.Dial("tcp", n.Addr+":"+string(n.Port))
		if er != nil {
			badNodeList = append(badNodeList, node.TNodeInfo{n.Addr, n.Port})
			// bad node
			continue
		}
		werr := FGWriter(c, fg)
		if werr != nil {
			badNodeList = append(badNodeList, node.TNodeInfo{n.Addr, n.Port})
			// write in error
			continue
		}
	}
	fBr.BadNodeList = badNodeList
	c <- *fBr
	return nil
}
