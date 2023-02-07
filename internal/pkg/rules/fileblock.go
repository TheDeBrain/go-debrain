package rules

// A set of rule constants used to constrain data partitioning within the FileBlock protocol.
// It must be consistent with the byte size of the data type in the FileBlock protocol.
const (
	// This is the size in bytes of the file index, size : 8 byte
	FILE_INDEX_DATASIZE_DESCRIPTOR_BYTE_NUM = 8
	// This is the size in bytes of the file name, size : 8 byte
	FILE_NAME_DATASIZE_BYTE_NUM = 8
	// Total packet size
	FILE_TOTAL_PACKET_SIZE_DESCRIPTOR_BYTE_NUM = 8
	// This is the size in bytes of the file block position, size : 4 byte
	FILE_BLOCK_POSITION_DATASIZE_DESCRIPTOR_BYTE_NUM = 4
	// This is the size in bytes of the file block size descriptor, size : 4 byte
	FILE_BLOCK_DATASIZE_DESCRIPTOR_BYTE_NUM = 4
	// This is the size in bytes of the time stamp, size : 8 byte
	FILE_BLOCK_TIMESTAMP_BYTE_NUM = 8
	// This is the size in bytes of the file owner, size : 4 byte
	FILE_BLOCK_OWNER_BYTE_NUM = 4
	// This is the size in bytes of the file block end flag, size : 4 byte
	FILE_BLOCK_END_FLAG_BYTE_NUM = 4
	// This is file block size constraint, size : 1024*1024*1 byte
	FILE_BLCOK_SIZE_CONSTRAINT = 1024 * 2
	// Thie is file block end flag
	FILE_BLCOK_END_FLAG = "end"
)
