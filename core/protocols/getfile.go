package protocols

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/derain/internal/pkg/rules"
	"net"
)

type GetFile struct {
	FileOwnerSize   uint32
	FileNameSize    uint64
	FileEndFlagSize uint32 // File block end flag data size , 4 bytz
	FileName        []byte // file name
	FileOwner       []byte // file owner data
	EndFlag         []byte // file block end flag
}

func GFNew(ows uint32, fns uint64, fn []byte, fo []byte) GetFile {
	gf := GetFile{
		ows,
		fns,
		uint32(len([]byte(rules.FILE_BLCOK_END_FLAG))),
		fn,
		fo,
		[]byte(rules.FILE_BLCOK_END_FLAG),
	}
	return gf
}

func GFBuf() {}

func GFNetUnPack(conn net.Conn) (GetFile, error) {
	gf := new(GetFile)
	// ---------------------------- protocol head ----------------------------
	// file owner size
	fileOwnerSizeBuf := make([]byte, rules.FILE_BLOCK_OWNER_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err := conn.Read(fileOwnerSizeBuf)
	fos := bytes.NewReader(fileOwnerSizeBuf)
	binary.Read(fos, binary.BigEndian, &gf.FileOwnerSize)
	if err != nil {
		return *gf, err
	}
	// file name size
	fileNameSizeBuf := make([]byte, rules.FILE_NAME_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileNameSizeBuf)
	fnsBuf := bytes.NewReader(fileNameSizeBuf)
	binary.Read(fnsBuf, binary.BigEndian, &gf.FileNameSize)
	if err != nil {
		return *gf, err
	}
	//  end size
	endBuf := make([]byte, rules.FILE_BLOCK_END_FLAG_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(endBuf)
	ebBuf := bytes.NewReader(endBuf)
	binary.Read(ebBuf, binary.BigEndian, &gf.FileEndFlagSize)
	if err != nil {
		return *gf, err
	}
	// ---------------------------- protocol body ----------------------------
	// file owner data
	fileOwner := make([]byte, gf.FileOwnerSize)
	fon, err := conn.Read(fileOwner)
	gf.FileOwner = fileOwner[:fon]
	if err != nil {
		return *gf, err
	}
	// file name data
	fileNameBuf := make([]byte, gf.FileNameSize)
	nn, err := conn.Read(fileNameBuf)
	gf.FileName = fileNameBuf[:nn]
	if err != nil {
		return *gf, err
	}
	// end data
	endB := make([]byte, gf.FileEndFlagSize)
	n, err := conn.Read(endB)
	gf.EndFlag = endB[:n]
	if err != nil {
		return *gf, err
	}
	if string(endB) == rules.FILE_BLCOK_END_FLAG {
		return *gf, nil
	}
	return *gf, errors.New("illegal agreement")
}
