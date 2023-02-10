package sync

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/derain/core/db/table/node"
	"github.com/derain/core/db/table/sys"
	"github.com/derain/core/protocols"
	"github.com/derain/internal/pkg/rules"
	"github.com/derain/internal/pkg/utils"
	"github.com/derain/test"
	"log"
	"net"
	"os"
	"strconv"
)

// start sync service
func StartSyncService() error {
	sys := sys.LoadTSys()
	ln, err := net.Listen("tcp", ":"+sys.SyncPort)
	defer ln.Close()
	if err != nil {
		return err
	}
	log.Println("Start the rpc server and listen to the port:", sys.SyncPort)
	for {
		conn, _ := ln.Accept()
		if err != nil {
			log.Fatal("sync connect error:", err)
		}
		go handleSyncService(conn)
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
		WFByFileToNet(fbuf, uint64(len(fbuf)), ptl)
	}
	return nil
}

// handle between server sync request
func handleBetweenServerSyncReq(conn net.Conn) error {
	// read in sync net
	err := RFByFileBlockToLocalInNet(conn)
	if err != nil {
		return err
	}
	// cync db
	return nil
}

// handle upload sync request
func handleUploadSyncReq(conn net.Conn) error {
	// read in sync net
	err := RFByFileBlockToLocalInNet(conn)
	if err != nil {
		return err
	}
	// cync db
	return nil
}

// handle send upload sync request
func HandleSendUploadSyncReq(b []byte) error {
	ptl := protocols.CommProtocol{
		ProtocolType: rules.FILE_BLOCK_UPLOAD_SYNC_PROTOCOL,
	}
	_, err := WFByFileToNet(b, uint64(len(b)), ptl)
	if err != nil {
		return err
	}
	return nil
}

// read file by FileBlock protocol in net
// Each time a FileBlock is successfully read, a pointer to the FileBlock will be returned
func RFByFileBlockToLocalInNet(conn net.Conn) error {
	// file block struct
	fb := new(protocols.FileBlock)
	// ---------------------------- protocol head ----------------------------
	// file index size
	fileIndexSizeBuf := make([]byte, rules.FILE_INDEX_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err := conn.Read(fileIndexSizeBuf)
	fisBuf := bytes.NewReader(fileIndexSizeBuf)
	binary.Read(fisBuf, binary.BigEndian, &fb.Head.FileIndexSize)
	if err != nil {
		return err
	}
	// file name size
	fileNameSizeBuf := make([]byte, rules.FILE_NAME_DATASIZE_BYTE_NUM)
	_, err = conn.Read(fileNameSizeBuf)
	fnsBuf := bytes.NewReader(fileNameSizeBuf)
	binary.Read(fnsBuf, binary.BigEndian, &fb.Head.FileNameSize)
	if err != nil {
		return err
	}
	// file total size
	fileTotalBuf := make([]byte, rules.FILE_TOTAL_SIZE_BYTE_NUM)
	_, err = conn.Read(fileTotalBuf)
	ftBuf := bytes.NewReader(fileTotalBuf)
	binary.Read(ftBuf, binary.BigEndian, &fb.Head.FileTotalSize)
	if err != nil {
		return err
	}
	// file total block num
	fileTotalBlockNumBuf := make([]byte, rules.FILE_TOTAL_BLOCK_NUM_BYTE_NUM)
	_, err = conn.Read(fileTotalBlockNumBuf)
	ftbBuf := bytes.NewReader(fileTotalBlockNumBuf)
	binary.Read(ftbBuf, binary.BigEndian, &fb.Head.FileTotalBlockNum)
	if err != nil {
		return err
	}
	// file block position
	fileBlockPositionBuf := make([]byte, rules.FILE_BLOCK_POSITION_DATASIZE_BYTE_NUM)
	_, err = conn.Read(fileBlockPositionBuf)
	fbpBuf := bytes.NewReader(fileBlockPositionBuf)
	binary.Read(fbpBuf, binary.BigEndian, &fb.Head.FileBlockPosition)
	if err != nil {
		return err
	}
	// file block size
	fileBlockSizeBuf := make([]byte, rules.FILE_BLOCK_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileBlockSizeBuf)
	bsbuf := bytes.NewReader(fileBlockSizeBuf)
	binary.Read(bsbuf, binary.BigEndian, &fb.Head.FileBlockSize)
	if err != nil {
		return err
	}
	// file owner size
	fileOwnerSizeBuf := make([]byte, rules.FILE_BLOCK_OWNER_BYTE_NUM)
	_, err = conn.Read(fileOwnerSizeBuf)
	fosbuf := bytes.NewReader(fileOwnerSizeBuf)
	binary.Read(fosbuf, binary.BigEndian, &fb.Head.FileOwnerSize)
	if err != nil {
		return err
	}
	// file storage node size
	fileBlockStorageNodeSizeBuf := make([]byte, rules.FILE_BLOCK_STROAGE_NODE_DATASIZE_BYTE_NUM)
	_, err = conn.Read(fileBlockStorageNodeSizeBuf)
	fsnbuf := bytes.NewReader(fileBlockStorageNodeSizeBuf)
	binary.Read(fsnbuf, binary.BigEndian, &fb.Head.FileBlockStorageNodeSize)
	if err != nil {
		return err
	}
	// file end flag size
	fileBlockEndFlagSizeBuf := make([]byte, rules.FILE_BLOCK_END_FLAG_BYTE_NUM)
	_, err = conn.Read(fileBlockEndFlagSizeBuf)
	febuf := bytes.NewReader(fileBlockEndFlagSizeBuf)
	binary.Read(febuf, binary.BigEndian, &fb.Head.FileBlockEndFlagSize)
	if err != nil {
		return err
	}
	// ---------------------------- protocol body ----------------------------
	// file index size
	fileIndexBuf := make([]byte, fb.Head.FileIndexSize)
	_, err = conn.Read(fileIndexBuf)
	fb.Body.FileIndex = fileIndexBuf
	if err != nil {
		return err
	}
	fileIndex := string(fileIndexBuf[:])
	fmt.Println(fileIndex)
	// file block storage node size
	fileStorageNodeBuf := make([]byte, fb.Head.FileBlockStorageNodeSize)
	_, err = conn.Read(fileStorageNodeBuf)
	fb.Body.FileBlockStorageNode = fileStorageNodeBuf
	if err != nil {
		return err
	}
	var s []node.TFBNodeInfo
	json.Unmarshal(fb.Body.FileBlockStorageNode[:], &s)
	// file name data
	fileNameBuf := make([]byte, fb.Head.FileNameSize)
	_, err = conn.Read(fileNameBuf)
	fb.Body.FileName = fileNameBuf
	if err != nil {
		return err
	}
	// file owner data
	fileOwnerBuf := make([]byte, fb.Head.FileOwnerSize)
	_, err = conn.Read(fileOwnerBuf)
	fb.Body.FileOwner = fileOwnerBuf
	if err != nil {
		return err
	}
	// file block data
	fileDataBuf := make([]byte, fb.Head.FileBlockSize)
	_, err = conn.Read(fileDataBuf)
	fb.Body.FileBlockData = fileDataBuf
	if err != nil {
		return err
	}
	fileName := string(fb.Body.FileName[:])
	fileOwner := string(fb.Body.FileOwner[:])
	// ---------------------------- protocol foot ----------------------------
	// file block end flag data
	fileBlockEndBuf := make([]byte, fb.Head.FileBlockEndFlagSize)
	n, err := conn.Read(fileBlockEndBuf)
	fb.Foot.FileBlockEndFlag = fileBlockEndBuf[:n]
	if err != nil {
		return err
	}
	endFlag := string(fb.Foot.FileBlockEndFlag)
	if endFlag == rules.FILE_BLCOK_END_FLAG {
		// file write
		fsys := sys.LoadFileSys()
		dir := fsys.FileStoragePath
		dir = dir + fileOwner + "/"
		if !utils.CheckPathExists(dir) {
			os.Mkdir(dir, 0777)
		}
		fileName = fileName + "_" + strconv.FormatInt(int64(fb.Head.FileBlockPosition), 10)
		//file storage
		fn := dir + fileName
		bb, _ := json.Marshal(fb)
		utils.WFToLocal(bb, fn)
		return nil
	}
	return nil
}

// write file to network by fileblock protocol
func WFByFileToNet(file []byte, fileSize uint64, ptl protocols.CommProtocol) (bool, error) {
	fr := bytes.NewReader(file)
	fbuf := make([]byte, len(file))
	fr.Read(fbuf)
	// file name
	fileName := test.TFName
	// file account
	ownerAddr := test.TOwner
	// route table
	routeTable := node.LoadRouteTable().NodeList
	// split file
	fl := utils.SplitFile(fbuf)
	// file position,increment
	FBPosition := uint32(0)
	// file uuid
	fUUID := utils.CrtUUID()
	// create channel list
	fBChannel := make(chan protocols.FileBlockSyncResult)
	// ---------------- write processing start ----------------
	for e := fl.Front(); e != nil; e = e.Next() {
		// create file block storage route table
		fBSNodeRoutable := new(node.TFBRouteTable)
		// create file blcok storage node list
		for _, nd := range routeTable {
			ip := nd.Addr
			port := nd.Port
			// node health check
			_, er := net.Dial("tcp", ip+":"+port)
			// bad node
			if er != nil {
				continue
			}
			// file block storage healthy node
			node := node.TFBNodeInfo{ip, port, fUUID, FBPosition}
			fBSNodeRoutable.NodeList = append(fBSNodeRoutable.NodeList, node)
		}
		// write in file block
		go wFByFileBlockToRouteTable(fUUID, fileName, fileSize, uint64(fl.Len()),
			FBPosition, e.Value.([]byte), ownerAddr, fBSNodeRoutable, ptl, fBChannel)
		// file block position increase
		FBPosition += 1
	}
	// ---------------- write processing end ----------------
	// file block sync result
	for i := 0; i < fl.Len(); i++ {
		res := <-fBChannel
		log.Println("file block sync result:", res)
	}
	// error hanlde
	return true, nil
}

// write file block to route table by fileblock protocol
func wFByFileBlockToRouteTable(fUUID string,
	fileName string,
	fileSize uint64,
	fileTotalBlockNum uint64,
	fileBlockPosition uint32,
	fileBlock []byte,
	ownerAddr string,
	fBSNodeRoutable *node.TFBRouteTable,
	ptl protocols.CommProtocol,
	c chan protocols.FileBlockSyncResult) {
	// file block result
	fBr := new(protocols.FileBlockSyncResult)
	// bad node
	var badNodeList []any
	// encode node list
	nLb, _ := json.Marshal(fBSNodeRoutable.NodeList)
	// create file block
	fs := protocols.FBNew(fUUID, fileName, fileSize, fileTotalBlockNum, fileBlockPosition, uint32(len(fileBlock)),
		uint32(len([]byte(ownerAddr))), uint64(len(nLb)), ownerAddr, rules.FILE_BLCOK_END_FLAG, nLb, fileBlock)
	//file block buffer
	buff := bytes.NewBuffer([]byte{})
	// read in protocol type
	protocols.CPBuf(buff, ptl)
	// read in file block protocol
	protocols.FBBuf(buff, fs)
	// file block sync
	for _, n := range fBSNodeRoutable.NodeList {
		c, er := net.Dial("tcp", n.Addr+":"+n.Port)
		if er != nil {
			badNodeList = append(badNodeList, node.TNodeInfo{n.Addr, n.Port})
			// bad node
			continue
		}
		_, werr := c.Write(buff.Bytes())
		if werr != nil {
			badNodeList = append(badNodeList, node.TNodeInfo{n.Addr, n.Port})
			// write in error
			continue
		}
	}
	fBr.BadNodeList = badNodeList
	c <- *fBr
}
