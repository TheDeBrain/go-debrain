package protocols

// Protocol for file block head
type FileBlockHead struct {
	ProtocolType      uint8  // Declare the protocol type , 1 byte
	FileIndex         uint64 // The index of the file to which the file block belongs , 8 byte
	FileBlockPosition uint32 // The space occupied by the file block in the file needs to be concatenated and read sequentially, starting from 0 , 4 byte
	FileBlockSize     uint32 // File block data size , 4 byte
	FileOwnerSize     uint32 // File block data size , 4 byte
	TimeStamp         int64  // time stamp unit : ms, 8 byte
}

// Protocol for file block body
type FileBlockBody struct {
	FileOwner     []byte // file owner data
	FileBlockData []byte // file block data
}

type FileBlock struct {
	Head FileBlockHead // file block head
	Body FileBlockBody // file block body
}