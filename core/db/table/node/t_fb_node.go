package node

// Node information in file blocks
type TFBNodeInfo struct {
	Addr              string `json:"addr"`                //node address
	Port              string `json:"port"`                //node port
	FileIndex         string `json:"file_id"`             // Index of the owning file
	FileBlockPosition uint32  `json:"file_block_position"` // The sequence number of the block in the file "FileId" stored by this node
}
