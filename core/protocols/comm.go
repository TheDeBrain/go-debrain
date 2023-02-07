package protocols

// Global communication protocol
type CommProtocol struct {
	ProtocolType uint8 // Declare the protocol type , 1 byte
	Data         any   // Protocol Data Body
}
