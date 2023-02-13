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
	ProtocolType    uint8 // Declare the protocol type , 1 byte
	FileOwnerSize   uint32
	FileNameSize    uint64
	FileEndFlagSize uint32 // File block end flag data size , 4 bytz
	FileName        []byte // file name
	FileOwner       []byte // file owner data
	EndFlag         []byte // file block end flag
}

func GFNew(ows uint32, fns uint64, fn []byte, fo []byte) *FileGetter {
	gf := FileGetter{
		rules.FILE_GETTER_PROTOCOL,
		ows,
		fns,
		uint32(len([]byte(rules.FILE_BLCOK_END_FLAG))),
		fn,
		fo,
		[]byte(rules.FILE_BLCOK_END_FLAG),
	}
	return &gf
}

func FGBuf(gf *FileGetter) (*bytes.Buffer, error) {
	buff := bytes.NewBuffer([]byte{})
	// read in protocol type
	binary.Write(buff, binary.BigEndian, gf.ProtocolType)
	// read in file index
	binary.Write(buff, binary.BigEndian, gf.FileOwnerSize)
	binary.Write(buff, binary.BigEndian, gf.FileNameSize)
	binary.Write(buff, binary.BigEndian, gf.FileEndFlagSize)
	binary.Write(buff, binary.BigEndian, gf.FileName)
	binary.Write(buff, binary.BigEndian, gf.FileOwner)
	binary.Write(buff, binary.BigEndian, rules.FILE_BLCOK_END_FLAG)
	return buff, nil
}

// file getter net un package
func FGNetUnPack(conn net.Conn) (*FileGetter, error) {
	fb := new(FileGetter)
	fb, err := FGReader(conn)
	if err != nil {
		return nil, err
	}
	return fb, err
}

// file getter protocol reader
func FGReader(r io.Reader) (*FileGetter, error) {
	gf := new(FileGetter)
	// ---------------------------- protocol head ----------------------------
	// protocol type
	protocolTypeBuf := make([]byte, int(unsafe.Sizeof(gf.ProtocolType)))
	_, err := r.Read(protocolTypeBuf)
	ptBuf := bytes.NewReader(protocolTypeBuf)
	binary.Read(ptBuf, binary.BigEndian, &gf.ProtocolType)
	if err != nil {
		return gf, err
	}
	// file owner size
	fileOwnerSizeBuf := make([]byte, int(unsafe.Sizeof(gf.FileOwnerSize)))
	_, err = r.Read(fileOwnerSizeBuf)
	fos := bytes.NewReader(fileOwnerSizeBuf)
	binary.Read(fos, binary.BigEndian, &gf.FileOwnerSize)
	if err != nil {
		return gf, err
	}
	// file name size
	fileNameSizeBuf := make([]byte, int(unsafe.Sizeof(gf.FileNameSize)))
	_, err = r.Read(fileNameSizeBuf)
	fnsBuf := bytes.NewReader(fileNameSizeBuf)
	binary.Read(fnsBuf, binary.BigEndian, &gf.FileNameSize)
	if err != nil {
		return gf, err
	}
	//  end size
	endBuf := make([]byte, int(unsafe.Sizeof(gf.FileEndFlagSize)))
	_, err = r.Read(endBuf)
	ebBuf := bytes.NewReader(endBuf)
	binary.Read(ebBuf, binary.BigEndian, &gf.FileEndFlagSize)
	if err != nil {
		return gf, err
	}
	// ---------------------------- protocol body ----------------------------
	// file owner data
	fileOwner := make([]byte, gf.FileOwnerSize)
	fon, err := r.Read(fileOwner)
	gf.FileOwner = fileOwner[:fon]
	if err != nil {
		return gf, err
	}
	// file name data
	fileNameBuf := make([]byte, gf.FileNameSize)
	nn, err := r.Read(fileNameBuf)
	gf.FileName = fileNameBuf[:nn]
	if err != nil {
		return gf, err
	}
	// end data
	endB := make([]byte, gf.FileEndFlagSize)
	n, err := r.Read(endB)
	gf.EndFlag = endB[:n]
	if err != nil {
		return gf, err
	}
	if string(endB) == rules.FILE_BLCOK_END_FLAG {
		return gf, nil
	}
	return gf, errors.New("illegal agreement")
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
	for _, n := range rt.NodeList {
		c, er := net.Dial("tcp", n.Addr+":"+n.Port)
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
