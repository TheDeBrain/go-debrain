package structs

// Protocol for file block head
type DBTableHead struct {
	ProtocolType uint8  // Declare the protocol type , 1 byte
	TableSize    uint64 // The , 8 byte
}

// Protocol for file block body
type DBTableBody struct {
	TableData []byte // file owner data
}

type DBTable struct {
	Head DBTableHead // file block head
	Body DBTableBody // file block body
}