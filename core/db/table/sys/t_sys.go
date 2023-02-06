package sys

// system table
type TSys struct {
	SyncPort          string `json:"sync_port"`          // sync
	RouteTablPath     string `json:"route_tabl_path"`    // route table path
	FileStoragePath   string `json:"file_storage_path"`  // file storage path
	HeartbeatInterval int    `json:"heartbeat_interval"` // node heartbeat detection interval, unit: second (s)
}
