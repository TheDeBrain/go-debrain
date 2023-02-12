package rules

const (
	// Full Mode.
	// Full synchronization will be performed.
	// The node with the most complete storage will send all the files in its local database to the synchronization requester
	FULL_MODE = 0 << (10 * iota)
	// Upload Mode.
	// In this mode, the network will process file upload. Generally, the number of file uploads will not be large, which will be far less than Full Mode
	UPLOAD_MODE
	// Download Mode
	// The network will send the required files to the requester according to the information of the requester
	DOWNLOAD_MODE
)
