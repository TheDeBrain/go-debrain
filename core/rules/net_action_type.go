package rules

//net action type
const (
	NET_ACTION_TYPE_SIZE = 1 // 1 byte
	FILE_BLOCK_CLIENT_SYNC_REQ = uint8(iota)
	FILE_BLOCK_CLIENT_SYNC_RECEIVE
	FILE_BLOCK_SERVER_BROADCAST_SYNC
	FILE_GETTER_RESPONSE

// server synchronization request,cata synchronization between service nodes
	FILE_BLOCK_BETWEEN_SERVER_SYNC_PROTOCOL
	// file upload synchronization request,file blocks are synchronized between network nodes for the first time
	FILE_BLOCK_UPLOAD_SYNC_PROTOCOL
	SYS_DB_SYNC_PROTOCOL
	FILE_SYS_DB_SYNC_PROTOCOL
)
