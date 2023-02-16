package sync

import (
	"container/list"
	"errors"
	"github.com/derain/core/db/table/node"
	"github.com/derain/core/db/table/sys"
	"github.com/derain/core/protocols"
	"github.com/derain/core/rules"
	"github.com/derain/test"
	"log"
	"net"
	"os"
	"path/filepath"
)

// handle send upload sync request
// b : file byte array
func HandleSendUploadSyncReq(file []byte, fileName string, fileOwner string) error {
	fbArr := protocols.FBNewArrayByFile(file, fileName, fileOwner)
	for _, fb := range fbArr {
		protocols.FBSyncFull(fb, rules.FILE_BLOCK_UPLOAD_SYNC_PROTOCOL)
	}
	return nil
}

// send get file block request
func HandleGetFileBlockReq(fileOwner string, fileName string) (*protocols.ResultCollect, error) {
	fob := []byte(fileOwner)
	fn := []byte(fileName)
	fg := protocols.FGNew(uint32(len(fob)), uint64(len(fn)), fob, fn)
	fgp, _ := protocols.FGBuf(fg)
	np := protocols.NPNew(rules.FILE_GETTER_PROTOCOL, fgp.Bytes())
	// route table
	rtr := node.RandomNodeGetter(rules.RANDOM_SYNC_NODE_NUM)
	rc, err := np.NPSendFull(rtr)
	if err != nil {
		return nil, err
	}
	return rc, nil
}

// Handle fetch file block requests
// Output file chunks to client
func handleGetFileResponse(conn net.Conn) error {
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
			res, _ := protocols.RESNew(conn.LocalAddr().String(), sys.TSysNew().SyncPort,
				rules.NET_PACK_OK_FLAG, "ok", fB)
			protocols.RESWriter(conn, res)
		}
	} else {
		// addr string, port string, flag string, des string
		res, _ := protocols.RESNew(conn.LocalAddr().String(), sys.TSysNew().SyncPort,
			rules.NET_PACK_ERROR_FLAG, "the file does not exist for this node", nil)
		protocols.RESWriter(conn, res)
		// file does not exist
		return errors.New("the file does not exist for this node")
	}
	return nil
}

// handle client sync request
func handleClientSyncReq(conn net.Conn) error {
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
		np := protocols.NPNew(rules.FILE_BLOCK_CLIENT_SYNC_RECEIVE, fbuf)
		protocols.NPWriter(conn, np)
	}
	return nil
}

// handle client sync request receive
func handleClientSyncReqReceive(con net.Conn) error {
	fb, err := protocols.FBReader(con)
	if err != nil {
		log.Println("handle client sync request receive error")
		return err
	}
	protocols.FBSaveToLocal(fb)
	return nil
}

// handle client sync request receive
func handleFileBlockServerBroadCastSync(con net.Conn) error {
	fb, err := protocols.FBReader(con)
	if err != nil {
		log.Println("handle client sync request receive error")
		return err
	}
	protocols.FBSaveToLocal(fb)
	return nil
}

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

// handle upload sync request
func handleUploadSyncReq(conn net.Conn) error {
	// read in sync net
	fb, err := protocols.FBNetUnPack(conn)
	if err != nil {
		return err
	}
	// save local
	protocols.FBSaveToLocal(fb)
	// result
	res, _ := protocols.RESNew(conn.LocalAddr().String(), sys.TSysNew().SyncPort,
		rules.NET_PACK_OK_FLAG, "ok", fb)
	protocols.RESWriter(conn, res)
	return nil
}
