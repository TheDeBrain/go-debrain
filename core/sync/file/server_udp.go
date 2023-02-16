package file

import (
	"errors"
	"github.com/derain/core/db/table/sys"
	"github.com/derain/core/protocols"
	"github.com/derain/core/rules"
	"log"
	"net"
	"path/filepath"
)

// Handle fetch file block requests
// Output file chunks to client
func HandleGetFileResponseUDP(conn *net.UDPConn) error {
	// get file protocol
	fg, err := protocols.FGNetUnPack(conn)
	if err != nil {
		return err
	}
	fileOwner := string(fg.FileOwner)
	fileName := string(fg.FileName)
	//file system
	fSys := sys.LoadFileSys()
	// storage path
	storagePath := fSys.FileStoragePath
	// file path
	filePath := storagePath
	filePath = filePath + fileOwner
	filePath = filePath + "/" + fileName
	m, err := filepath.Glob(filePath + "_" + "[1-3]")
	if err != nil {
		return err
	}
	// the one file
	if len(m) > 0 {
		var fBs []*protocols.FileBlock
		for _, f := range m {
			val := f
			fB, err := protocols.RFBByPath(val)
			if err != nil {
				return err
			}
			// file block node route
			//nodeList := new([]node.TFBNodeInfo)
			// decode
			//json.Unmarshal(fB.Body.FileBlockStorageNode, &nodeList)
			//for _, node := range *nodeList {
			//	fmt.Println("node list--", node)
			//}
			fBs = append(fBs, fB)
		}
		res, _ := protocols.RESNew(conn.LocalAddr().String(), string(sys.TSysNew().SyncPortUDP),
			rules.NET_PACK_OK_FLAG, "ok", fBs)
		protocols.RESWriter(conn, res)
	} else {
		// addr string, port string, flag string, des string
		res, _ := protocols.RESNew(conn.LocalAddr().String(), string(sys.TSysNew().SyncPortUDP),
			rules.NET_PACK_ERROR_FLAG, "the file does not exist for this node", nil)
		protocols.RESWriter(conn, res)
		// file does not exist
		return errors.New("the file does not exist for this node")
	}
	return nil
}

// handle upload sync request
func HandleUploadSyncReqUDP(conn *net.UDPConn) error {
	// read in sync net
	fb, err := protocols.FBNetUnPack(conn)
	if err != nil {
		return err
	}
	// save local
	protocols.FBSaveToLocal(fb)
	// result
	res, _ := protocols.RESNew(conn.LocalAddr().String(), string(sys.TSysNew().SyncPortUDP),
		rules.NET_PACK_OK_FLAG, "ok", []*protocols.FileBlock{fb})
	protocols.RESWriter(conn, res)
	return nil
}

// handle file block server broad cast sync
func HandleFileBlockServerBroadCastSyncUDP(conn *net.UDPConn) error {
	fb, err := protocols.FBReader(conn)
	if err != nil {
		log.Println("handle client sync request receive error")
		return err
	}
	protocols.FBSaveToLocal(fb)
	return nil
}
