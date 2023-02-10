package protocols

import (
	"bytes"
)

// Protocol for database head
type DBTableHead struct {
	TablSize uint64 // db table size
}

// Protocol for database body
type DBTableBody struct {
	Table []byte // file owner data
}

// Protocol for database foot
type DBTableFoot struct {
	EndFlag []byte // file block end flag
}

// Protocol for database
type DBTable struct {
	Head DBTableHead // file block head
	Body DBTableBody // file block body
	Foot DBTableFoot // file block foot
}

func DBTNew(fileIndex string, fileName string, fileSize uint64,
	fileTotalBlockNum uint64, fileBlockPosition uint32,
	fileBlockSize uint32, fileOwnerSize uint32, fileBlockStorageNodeSize uint64,
	ownerAddr string,
	fileBlockEndFlag string,
	nLb []byte,
	fileBlockData []byte) DBTable {
}

func DBTBuf(buff *bytes.Buffer, dbt DBTable) {}
