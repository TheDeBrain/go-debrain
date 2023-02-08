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
	FileIndex                uint64 // The index of the file to which the file block belongs , 8 byte
	FileNameSize             uint64 // The file name of the file to which the file block belongs , 8 byte
	FileTotalSize            uint64 // The total size of the file to which the file block belongs , 8 byte
	FileBlockPosition        uint32 // The space occupied by the file block in the file needs to be concatenated and read sequentially, starting from 0 , 4 byte
	FileBlockSize            uint32 // File block data size , 4 byte
	FileOwnerSize            uint32 // File block data size , 4 byte
	FileBlockStorageNodeSize uint64 // File block storage node , 8 byte
	FileBlockEndFlagSize     uint32 // File block data size , 4 byte
}

// Protocol for file block body
type FileBlockBody struct {
	FileBlockStorageNode []byte // file block storage node
	FileName             []byte // file name
	FileOwner            []byte // file owner data
	FileBlockData        []byte // file block data
}

// Protocol for file block foot
type FileBlockFoot struct {
	FileBlockEndFlag []byte // file block end flag
}

func FBBuf(buff *bytes.Buffer, fs FileBlock) {
	// read in file index
	binary.Write(buff, binary.BigEndian, fs.Head.FileIndex)
	// read in file name
	binary.Write(buff, binary.BigEndian, fs.Head.FileNameSize)
	// read in file total size,unit:byte
	binary.Write(buff, binary.BigEndian, fs.Head.FileTotalSize)
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
