package protocols

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/derain/core/db/table/node"
	"github.com/derain/internal/pkg/rules"
	"net"
)

// Protocol for file block
type FileBlock struct {
	Head FileBlockHead // file block head
	Body FileBlockBody // file block body
	Foot FileBlockFoot // file block foot
}

// Protocol for file block head
type FileBlockHead struct {
	FileIndexSize            uint64 // The index of the file to which the file block belongs , 8 byte
	FileNameSize             uint64 // The file name of the file to which the file block belongs , 8 byte
	FileTotalSize            uint64 // The total size of the file to which the file block belongs , 8 byte
	FileTotalBlockNum        uint64 // The total block num of the file to which the file block belongs , 8 byte
	FileBlockPosition        uint32 // The space occupied by the file block in the file needs to be concatenated and read sequentially, starting from 0 , 4 byte
	FileBlockSize            uint32 // File block data size , 4 byte
	FileOwnerSize            uint32 // File block data size , 4 byte
	FileBlockStorageNodeSize uint64 // File block storage node , 8 byte
	FileBlockEndFlagSize     uint32 // File block end flag data size , 4 byte
}

// Protocol for file block body
type FileBlockBody struct {
	FileIndex            []byte // file index
	FileBlockStorageNode []byte // file block storage node
	FileName             []byte // file name
	FileOwner            []byte // file owner data
	FileBlockData        []byte // file block data
}

// Protocol for file block foot
type FileBlockFoot struct {
	EndFlag []byte // file block end flag
}

// Protocol Result
type FileBlockSyncResult struct {
	BadNodeList []any // bad node list
}

func FBNew(fileIndex string, fileName string, fileSize uint64,
	fileTotalBlockNum uint64, fileBlockPosition uint32,
	fileBlockSize uint32, fileOwnerSize uint32, fileBlockStorageNodeSize uint64,
	ownerAddr string,
	fileBlockEndFlag string,
	nLb []byte,
	fileBlockData []byte) FileBlock {
	fs := FileBlock{
		Head: FileBlockHead{
			FileIndexSize:            uint64(len([]byte(fileIndex))), // uuid
			FileNameSize:             uint64(len([]byte(fileName))),
			FileTotalSize:            fileSize,
			FileTotalBlockNum:        fileTotalBlockNum,
			FileBlockPosition:        fileBlockPosition,
			FileBlockSize:            fileBlockSize,
			FileOwnerSize:            fileOwnerSize,
			FileBlockStorageNodeSize: fileBlockStorageNodeSize,
			FileBlockEndFlagSize:     uint32(len([]byte(fileBlockEndFlag))),
		},
		Body: FileBlockBody{
			FileIndex:            []byte(fileIndex),
			FileBlockStorageNode: nLb,
			FileName:             []byte(fileName),
			FileOwner:            []byte(ownerAddr),
			FileBlockData:        fileBlockData,
		},
		Foot: FileBlockFoot{
			EndFlag: []byte(fileBlockEndFlag),
		},
	}
	return fs
}

func FBBuf(buff *bytes.Buffer, fs FileBlock) {
	// read in file index
	binary.Write(buff, binary.BigEndian, fs.Head.FileIndexSize)
	// read in file name
	binary.Write(buff, binary.BigEndian, fs.Head.FileNameSize)
	// read in file total size,unit:byte
	binary.Write(buff, binary.BigEndian, fs.Head.FileTotalSize)
	// read in file total block num
	binary.Write(buff, binary.BigEndian, fs.Head.FileTotalBlockNum)
	// read in file block position
	binary.Write(buff, binary.BigEndian, fs.Head.FileBlockPosition)
	// read in file block size
	binary.Write(buff, binary.BigEndian, fs.Head.FileBlockSize)
	// read in file owner size
	binary.Write(buff, binary.BigEndian, fs.Head.FileOwnerSize)
	// read in file block storage node size
	binary.Write(buff, binary.BigEndian, fs.Head.FileBlockStorageNodeSize)
	// read in file end flag size
	binary.Write(buff, binary.BigEndian, fs.Head.FileBlockEndFlagSize)
	// read in file index
	binary.Write(buff, binary.BigEndian, fs.Body.FileIndex)
	// read in file block storage node
	binary.Write(buff, binary.BigEndian, fs.Body.FileBlockStorageNode)
	// read in file name
	binary.Write(buff, binary.BigEndian, fs.Body.FileName)
	// read in file owner
	binary.Write(buff, binary.BigEndian, fs.Body.FileOwner)
	// read in file block data
	binary.Write(buff, binary.BigEndian, fs.Body.FileBlockData)
	// read in file block end flag
	binary.Write(buff, binary.BigEndian, fs.Foot.EndFlag)
}

// network unpack by fileblock
func FBNetUnPack(conn net.Conn) (FileBlock, error) {
	fb := new(FileBlock)
	// ---------------------------- protocol head ----------------------------
	// file index size
	fileIndexSizeBuf := make([]byte, rules.FILE_INDEX_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err := conn.Read(fileIndexSizeBuf)
	fisBuf := bytes.NewReader(fileIndexSizeBuf)
	binary.Read(fisBuf, binary.BigEndian, &fb.Head.FileIndexSize)
	if err != nil {
		return *fb, err
	}
	// file name size
	fileNameSizeBuf := make([]byte, rules.FILE_NAME_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileNameSizeBuf)
	fnsBuf := bytes.NewReader(fileNameSizeBuf)
	binary.Read(fnsBuf, binary.BigEndian, &fb.Head.FileNameSize)
	if err != nil {
		return *fb, err
	}
	// file total size
	fileTotalBuf := make([]byte, rules.FILE_TOTAL_SIZE_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileTotalBuf)
	ftBuf := bytes.NewReader(fileTotalBuf)
	binary.Read(ftBuf, binary.BigEndian, &fb.Head.FileTotalSize)
	if err != nil {
		return *fb, err
	}
	// file total block num
	fileTotalBlockNumBuf := make([]byte, rules.FILE_TOTAL_BLOCK_NUM_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileTotalBlockNumBuf)
	ftbBuf := bytes.NewReader(fileTotalBlockNumBuf)
	binary.Read(ftbBuf, binary.BigEndian, &fb.Head.FileTotalBlockNum)
	if err != nil {
		return *fb, err
	}
	// file block position
	fileBlockPositionBuf := make([]byte, rules.FILE_BLOCK_POSITION_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileBlockPositionBuf)
	fbpBuf := bytes.NewReader(fileBlockPositionBuf)
	binary.Read(fbpBuf, binary.BigEndian, &fb.Head.FileBlockPosition)
	if err != nil {
		return *fb, err
	}
	// file block size
	fileBlockSizeBuf := make([]byte, rules.FILE_BLOCK_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileBlockSizeBuf)
	bsbuf := bytes.NewReader(fileBlockSizeBuf)
	binary.Read(bsbuf, binary.BigEndian, &fb.Head.FileBlockSize)
	if err != nil {
		return *fb, err
	}
	// file owner size
	fileOwnerSizeBuf := make([]byte, rules.FILE_BLOCK_OWNER_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileOwnerSizeBuf)
	fosbuf := bytes.NewReader(fileOwnerSizeBuf)
	binary.Read(fosbuf, binary.BigEndian, &fb.Head.FileOwnerSize)
	if err != nil {
		return *fb, err
	}
	// file storage node size
	fileBlockStorageNodeSizeBuf := make([]byte, rules.FILE_BLOCK_STROAGE_NODE_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileBlockStorageNodeSizeBuf)
	fsnbuf := bytes.NewReader(fileBlockStorageNodeSizeBuf)
	binary.Read(fsnbuf, binary.BigEndian, &fb.Head.FileBlockStorageNodeSize)
	if err != nil {
		return *fb, err
	}
	// file end flag size
	fileBlockEndFlagSizeBuf := make([]byte, rules.FILE_BLOCK_END_FLAG_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileBlockEndFlagSizeBuf)
	febuf := bytes.NewReader(fileBlockEndFlagSizeBuf)
	binary.Read(febuf, binary.BigEndian, &fb.Head.FileBlockEndFlagSize)
	if err != nil {
		return *fb, err
	}
	// ---------------------------- protocol body ----------------------------
	// file index size
	fileIndexBuf := make([]byte, fb.Head.FileIndexSize)
	_, err = conn.Read(fileIndexBuf)
	fb.Body.FileIndex = fileIndexBuf
	if err != nil {
		return *fb, err
	}
	fileIndex := string(fileIndexBuf[:])
	fmt.Println(fileIndex)
	// file block storage node size
	fileStorageNodeBuf := make([]byte, fb.Head.FileBlockStorageNodeSize)
	_, err = conn.Read(fileStorageNodeBuf)
	fb.Body.FileBlockStorageNode = fileStorageNodeBuf
	if err != nil {
		return *fb, err
	}
	var s []node.TFBNodeInfo
	json.Unmarshal(fb.Body.FileBlockStorageNode[:], &s)
	// file name data
	fileNameBuf := make([]byte, fb.Head.FileNameSize)
	_, err = conn.Read(fileNameBuf)
	fb.Body.FileName = fileNameBuf
	if err != nil {
		return *fb, err
	}
	// file owner data
	fileOwnerBuf := make([]byte, fb.Head.FileOwnerSize)
	_, err = conn.Read(fileOwnerBuf)
	fb.Body.FileOwner = fileOwnerBuf
	if err != nil {
		return *fb, err
	}
	// file block data
	fileDataBuf := make([]byte, fb.Head.FileBlockSize)
	_, err = conn.Read(fileDataBuf)
	fb.Body.FileBlockData = fileDataBuf
	if err != nil {
		return *fb, err
	}
	// ---------------------------- protocol foot ----------------------------
	// file block end flag data
	fileBlockEndBuf := make([]byte, fb.Head.FileBlockEndFlagSize)
	n, err := conn.Read(fileBlockEndBuf)
	fb.Foot.EndFlag = fileBlockEndBuf[:n]
	if err != nil {
		return *fb, err
	}
	endFlag := string(fb.Foot.EndFlag)
	if endFlag == rules.FILE_BLCOK_END_FLAG {
		return *fb, nil
	}
	return *fb, errors.New("illegal agreement")
}
