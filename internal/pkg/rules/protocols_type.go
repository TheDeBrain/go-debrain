package rules

// protocol type identifier
const (
	// protocol type
	PROTOCOL_TYPE_BYTE_NUM = 1
	FILE_BLOCK_SYNC_PROTOCOL = 0 << (10 * iota)
	SYS_DB_SYNC_PROTOCOL
	FILE_SYS_DB_SYNC_PROTOCOL
)
