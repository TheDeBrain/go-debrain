package sync

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/derain/core/db/table/node"
	"github.com/derain/core/db/table/sys"
	"github.com/derain/core/protocols"
	"github.com/derain/core/rules"
	"github.com/derain/test"
	"net"
	"os"
	"path/filepath"
)

// handle send upload sync request
func HandleSendUploadSyncReq(b []byte) error {
	ptl := protocols.CommProtocol{
		ProtocolType: rules.FILE_BLOCK_UPLOAD_SYNC_PROTOCOL,
	}
	_, err := protocols.WFBToNet(b, uint64(len(b)), &ptl)
	if err != nil {
		return err
	}
	return nil
}

// send get file block request
func HandleGetFileBlockReq(fileOwner string, fileName string) error {
	fob := []byte(fileOwner)
	fn := []byte(fileName)
	  protocols.GFNew(uint32(len(fob)), uint64(len(fn)), fob, fn)

	return nil
}

// Handle fetch file block requests
// Output file chunks to client
func handleGetFileResponse(conn net.Conn) error {
	// get file protocol
	gf, err := protocols.GFNetUnPack(conn)
	if err != nil {
		return err
	}
	fileOwner := string(gf.FileOwner)
	fileName := string(gf.FileName)
	//file system
	fSys := sys.LoadFileSys()
	// storage path
	storagePath := fSys.FileStoragePath
	// file path
	filePath := storagePath
	filePath = filePath + fileOwner
	filePath = filePath + "/" + fileName
	m, err := filepath.Glob(filePath + "[1-2]")
	if err != nil {
		return err
	}
	// the one file
	if len(m) > 0 {
		val := m[0]
		fB, err := protocols.RFBInLocal(val)
		if err != nil {
			return err
		}
		// file block node route
		nodeList := new([]node.TFBNodeInfo)
		json.Unmarshal(fB.Body.FileBlockStorageNode, &nodeList)
		fmt.Println("文件块节点路由")
		fmt.Println(nodeList)
		fmt.Println(val)
	} else {
		// file does not exist
		return errors.New("the file does not exist for this node")
	}
	return nil
}

// handle client sync request
func handleClientSyncReq(conn net.Conn) error {
	ptl := protocols.CommProtocol{
		ProtocolType: rules.FILE_BLOCK_UPLOAD_SYNC_PROTOCOL,
	}
	fileList := list.New()
	// test
	fileList.PushBack(test.Path1)
	for e := fileList.Front(); e != nil; e = e.Next() {
		filePath := e.Value.(string)
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			return err
		}
		fInfo, e := os.Stat(filePath)
		if e != nil {
			return e
		}
		fbuf := make([]byte, fInfo.Size())
		f.Read(fbuf)
		// write in sync net
		protocols.WFBToNet(fbuf, uint64(len(fbuf)), &ptl)
	}
	return nil
}

// handle between server sync request
func handleBetweenServerSyncReq(conn net.Conn) error {
	// read in sync net
	err := protocols.RFBToLocalInNet(conn)
	if err != nil {
		return err
	}
	// cync db
	return nil
}

// handle upload sync request
func handleUploadSyncReq(conn net.Conn) error {
	// read in sync net
	err := protocols.RFBToLocalInNet(conn)
	if err != nil {
		return err
	}
	// cync db
	return nil
}
