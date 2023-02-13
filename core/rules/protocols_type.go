package rules

// protocol type
const (
	// file block protocol
	FILE_BLOCK_PROTOCOL = uint8(1 >> iota)
	// file getter protocol
	FILE_GETTER_PROTOCOL
	// result collect protocol
	RESULT_COLLECT_PROTOCOL
	// result protocol
	RESULT_PROTOCOL
)
