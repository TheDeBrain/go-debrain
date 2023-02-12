package rules

// A set of rule constants used to constrain data partitioning within the FileBlock protocol.
// It must be consistent with the byte size of the data type in the FileBlock protocol.
const (
	// This is file block size constraint, size : 1024*1024*1 byte
	FILE_BLCOK_SIZE_CONSTRAINT = 1024 * 2
	// Network single maximum file processing capacity
	MAX_FILE_SIZE = 1024 * 1024 * 50
	// Thie is file block end flag
	FILE_BLCOK_END_FLAG = "end"
)
