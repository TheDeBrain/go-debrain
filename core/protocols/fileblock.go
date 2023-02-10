package protocols

import (
	"bytes"
	"encoding/binary"
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
	FileBlockEndFlagSize     uint32 // File block data size , 4 byte
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
	FileBlockEndFlag []byte // file block end flag
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
			FileBlockEndFlag: []byte(fileBlockEndFlag),
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
	binary.Write(buff, binary.BigEndian, fs.Foot.FileBlockEndFlag)
}
