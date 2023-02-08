package sync

// Protocol for file block head
type FileBlockHead struct {
	FileIndex            uint64 // The index of the file to which the file block belongs , 8 byte
	FileNameSize         uint64 // The file name of the file to which the file block belongs , 8 byte
	FileTotalSize        uint64 // The total size of the file to which the file block belongs , 8 byte
	FileBlockPosition    uint32 // The space occupied by the file block in the file needs to be concatenated and read sequentially, starting from 0 , 4 byte
	FileBlockSize        uint32 // File block data size , 4 byte
	FileOwnerSize        uint32 // File block data size , 4 byte
	FileBlockEndFlagSize uint32 // File block data size , 4 byte
}

// Protocol for file block body
type FileBlockBody struct {
	FileName      []byte // file name
	FileOwner     []byte // file owner data
	FileBlockData []byte // file block data
}

// Protocol for file block foot
type FileBlockFoot struct {
	FileBlockEndFlag []byte // file block end flag
}

// Protocol for file block
type FileBlock struct {
	Head FileBlockHead // file block head
	Body FileBlockBody // file block body
	Foot FileBlockFoot // file block foot
}
