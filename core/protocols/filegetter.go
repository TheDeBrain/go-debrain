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

type FileGetter struct {
	FileOwnerSize   uint32
	FileNameSize    uint64
	FileEndFlagSize uint32 // File block end flag data size , 4 bytz
	FileName        []byte // file name
	FileOwner       []byte // file owner data
	EndFlag         []byte // file block end flag
}

func GFNew(ows uint32, fns uint64, fn []byte, fo []byte) *FileGetter {
	gf := FileGetter{
		ows,
		fns,
		uint32(len([]byte(rules.FILE_BLCOK_END_FLAG))),
		fn,
		fo,
		[]byte(rules.FILE_BLCOK_END_FLAG),
	}
	return &gf
}

func FGBuf(buff *bytes.Buffer, gf *FileGetter) {
	// read in file index
	binary.Write(buff, binary.BigEndian, gf.FileOwnerSize)
	binary.Write(buff, binary.BigEndian, gf.FileNameSize)
	binary.Write(buff, binary.BigEndian, gf.FileEndFlagSize)
	binary.Write(buff, binary.BigEndian, gf.FileName)
	binary.Write(buff, binary.BigEndian, gf.FileOwner)
	binary.Write(buff, binary.BigEndian, rules.FILE_BLCOK_END_FLAG)
}

// file getter net un package
func FGNetUnPack(conn net.Conn) (*FileGetter, error) {
	fb := new(FileGetter)
	fb, err := FGProtocolAnalysis(conn)
	if err != nil {
		return nil, err
	}
	return fb, err
}

// write file block to route table by fileblock protocol
func WFByGetFileToRouteTable(
	fileGetter *FileGetter,
	fBSNodeRoutable *node.TFBRouteTable,
	c chan FileBlockSyncResult) {
	// file block result
	fBr := new(FileBlockSyncResult)
	// bad node
	var badNodeList []any
	//file block buffer
	buff := bytes.NewBuffer([]byte{})
	// read in file block protocol
	FGBuf(buff, fileGetter)
	// file block sync
	for _, n := range fBSNodeRoutable.NodeList {
		c, er := net.Dial("tcp", n.Addr+":"+n.Port)
		if er != nil {
			badNodeList = append(badNodeList, node.TNodeInfo{n.Addr, n.Port})
			// bad node
			continue
		}
		_, werr := c.Write(buff.Bytes())
		if werr != nil {
			badNodeList = append(badNodeList, node.TNodeInfo{n.Addr, n.Port})
			// write in error
			continue
		}
	}
	fBr.BadNodeList = badNodeList
	c <- *fBr
}

// file getter protocol analysis in reader steam
func FGProtocolAnalysis(r io.Reader) (*FileGetter, error) {
	gf := new(FileGetter)
	// ---------------------------- protocol head ----------------------------
	// file owner size
	fileOwnerSizeBuf := make([]byte, int(unsafe.Sizeof(gf.FileOwnerSize)))
	_, err := r.Read(fileOwnerSizeBuf)
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
