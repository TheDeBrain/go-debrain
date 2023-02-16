package sync

import (
	"github.com/derain/core/protocols"
	"net"
)

// handle between server sync request
func handleBetweenServerSyncReq(conn net.Conn) error {
	// read in sync net
	fb, err := protocols.FBNetUnPack(conn)
	if err != nil {
		return err
	}
	// save local
	protocols.FBSaveToLocal(fb)
	return nil
}


