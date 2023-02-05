package rules

// A set of rule constants used to constrain data partitioning within the FileBlock protocol.
// It must be consistent with the byte size of the data type in the FileBlock protocol.
const (
	// This is the size in bytes of the file index
	FILE_INDEX_DATASIZE_DESCRIPTOR_BYTE_NUM = 8
	// This is the size in bytes of the file block position
	FILE_BLOCK_POSITION_DATASIZE_DESCRIPTOR_BYTE_NUM = 4
	// This is the size in bytes of the file block size descriptor
	FILE_BLOCK_DATASIZE_DESCRIPTOR_BYTE_NUM = 4
	// This is the size in bytes of the time stamp
	FILE_BLOCK_TIMESTAMP_BYTE_NUM = 8
	// This is Owner data size byte num
	FILE_OWNER_DATASIZE_BYTE_NUM = 100
	// This is protocol end flag byte num
	END_FLAG_BYTE_NUM = 1
)
