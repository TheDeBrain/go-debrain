package sync

// Protocol for database head
type DBTableHead struct {
	ProtocolType      uint8  // Declare the protocol type , 1 byte
	FileIndex         uint64 // The index of the file to which the file block belongs , 8 byte
	FileBlockPosition uint32 // The space occupied by the file block in the file needs to be concatenated and read sequentially, starting from 0 , 4 byte
	FileBlockSize     uint32 // File block data size , 4 byte
	FileOwnerSize     uint32 // File block data size , 4 byte
	FileBlockEndSize  uint32 // File block data size , 4 byte
	TimeStamp         int64  // time stamp unit : ms, 8 byte
}

// Protocol for database body
type DBTableBody struct {
	FileOwner     []byte // file owner data
	FileBlockData []byte // file block data
}

// Protocol for database foot
type DBTableFoot struct {
	FileBlockEndFlag []byte // file block end flag
}

// Protocol for database
type DBTable struct {
	Head DBTableHead // file block head
	Body DBTableBody // file block body
	Foot DBTableFoot // file block foot
}
