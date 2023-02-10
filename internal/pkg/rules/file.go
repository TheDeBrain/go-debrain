package rules

// A set of rule constants used to constrain data partitioning within the FileBlock protocol.
// It must be consistent with the byte size of the data type in the FileBlock protocol.
// Naming conventions
// DATASIZE DESCRIPTOR : Data Descriptor Byte Size
const (
	// This is the size in bytes of the file index, size : 8 byte
	FILE_INDEX_DATASIZE_DESCRIPTOR_BYTE_NUM = 8
	// This is the size in bytes of the file name, size : 8 byte
	FILE_NAME_DATASIZE_DESCRIPTOR_BYTE_NUM = 8
	// Total packet size byte num, size : 8 byte
	FILE_TOTAL_SIZE_DATASIZE_DESCRIPTOR_BYTE_NUM = 8
	// Total block num byte num, size : 8 byte
	FILE_TOTAL_BLOCK_NUM_DATASIZE_DESCRIPTOR_BYTE_NUM = 8
	// This is the size in bytes of the file block position, size : 4 byte
	FILE_BLOCK_POSITION_DATASIZE_DESCRIPTOR_BYTE_NUM = 4
	// This is the size in bytes of the file block size descriptor, size : 4 byte
	FILE_BLOCK_DATASIZE_DESCRIPTOR_BYTE_NUM = 4
	// This is the size in bytes of the time stamp, size : 8 byte
	FILE_BLOCK_TIMESTAMP_BYTE_NUM = 8
	// This is the size in bytes of the file stroage node byte num, size : 8 byte
	FILE_BLOCK_STROAGE_NODE_DATASIZE_DESCRIPTOR_BYTE_NUM = 8
	// This is the size in bytes of the file owner, size : 4 byte
	FILE_BLOCK_OWNER_DATASIZE_DESCRIPTOR_BYTE_NUM = 4
	// This is the size in bytes of the file block end flag, size : 4 byte
	FILE_BLOCK_END_FLAG_DATASIZE_DESCRIPTOR_BYTE_NUM = 4
	// This is file block size constraint, size : 1024*1024*1 byte
	FILE_BLCOK_SIZE_CONSTRAINT = 1024 * 2
	// Network single maximum file processing capacity
	MAX_FILE_SIZE = 1024 * 1024 * 50
	// Thie is file block end flag
	FILE_BLCOK_END_FLAG = "end"
)
