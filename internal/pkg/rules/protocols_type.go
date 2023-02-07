package rules

// protocol type identifier
const (
	// protocol type
	PROTOCOL_TYPE_BYTE_NUM = 1
	// client synchronization request,file block sync request,clients need to synchronize data to their local
	FILE_BLOCK_CLIENT_SYNC_PROTOCOL = iota
	// server synchronization request,cata synchronization between service nodes
	FILE_BLOCK_BETWEEN_SERVER_SYNC_PROTOCOL
	// file upload synchronization request,file blocks are synchronized between network nodes for the first time
	FILE_BLOCK_UPLOAD_SYNC_PROTOCOL
	SYS_DB_SYNC_PROTOCOL
	FILE_SYS_DB_SYNC_PROTOCOL
)
