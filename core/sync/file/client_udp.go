package file

import (
	"container/list"
	"github.com/derain/core/db/table/node"
	"github.com/derain/core/protocols"
	"github.com/derain/core/rules"
	"github.com/derain/test"
	"log"
	"net"
	"os"
)

// handle client sync request
func HandleClientSyncReqUDP(conn *net.UDPConn) error {
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
func HandleClientSyncReqReceiveUDP(conn *net.UDPConn) error {
	fb, err := protocols.FBReader(conn)
	if err != nil {
		log.Println("handle client sync request receive error")
		return err
	}
	protocols.FBSaveToLocal(fb)
	return nil
}

// send get file block request
func HandleGetFileBlockReqUDP(fileOwner string, fileName string) error {
	fob := []byte(fileOwner)
	fn := []byte(fileName)
	fg := protocols.FGNew(uint32(len(fob)), uint64(len(fn)), fob, fn)
	fgp, _ := protocols.FGBuf(fg)
	np := protocols.NPNew(rules.FILE_GETTER_PROTOCOL, fgp.Bytes())
	// route table
	rtr := node.RandomNodeGetter(rules.RANDOM_SYNC_NODE_NUM,"udp")
	err := np.NPSendFullUDP(rtr)
	if err != nil {
		return err
	}
	return nil
}

// handle send upload sync request
// b : file byte array
func HandleSendUploadSyncReqUDP(file []byte, fileName string, fileOwner string) error {
	fbArr := protocols.FBNewArrayByFile(file, fileName, fileOwner)
	for _, fb := range fbArr {
		protocols.FBSyncFull(fb, rules.FILE_BLOCK_UPLOAD_SYNC_PROTOCOL, "udp")
	}
	return nil
}
