package protocols

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/derain/core/db/table/node"
	"github.com/derain/core/db/table/sys"
	"github.com/derain/core/rules"
	"github.com/derain/internal/pkg/utils"
	"github.com/derain/test"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"unsafe"
)

// Protocol for file block
type FileBlock struct {
	Head FileBlockHead // file block head
	Body FileBlockBody // file block body
	Foot FileBlockFoot // file block foot
}

// Protocol for file block head
type FileBlockHead struct {
	ProtocolType             uint8  // Declare the protocol type , 1 byte
	FileIndexSize            uint64 // The index of the file to which the file block belongs , 8 byte
	FileNameSize             uint64 // The file name of the file to which the file block belongs , 8 byte
	FileTotalSize            uint64 // The total size of the file to which the file block belongs , 8 byte
	FileTotalBlockNum        uint64 // The total block num of the file to which the file block belongs , 8 byte
	FileBlockPosition        uint32 // The space occupied by the file block in the file needs to be concatenated and read sequentially, starting from 0 , 4 byte
	FileBlockSize            uint32 // File block data size , 4 byte
	FileOwnerSize            uint32 // File block data size , 4 byte
	FileBlockStorageNodeSize uint64 // File block storage node , 8 byte
	FileBlockEndFlagSize     uint32 // File block end flag data size , 4 byte
}

// Protocol for file block body
type FileBlockBody struct {
	FileIndex            []byte // file index
	FileBlockStorageNode []byte // file block storage node
	FileName             []byte // file name
	FileOwner            []byte // file owner data
	FileBlockData        []byte // file block data
}

// Protocol for file block foot
type FileBlockFoot struct {
	EndFlag []byte // file block end flag
}

// Protocol Result
type FileBlockSyncResult struct {
	BadNodeList []any // bad node list
}

// ------------------------- struct handle start -------------------------

// create file block pointer
func FBNew(fileIndex string, fileName string, fileSize uint64,
	fileTotalBlockNum uint64, fileBlockPosition uint32,
	fileBlockSize uint32, fileOwnerSize uint32, fileBlockStorageNodeSize uint64,
	fileOwner string,
	fileBlockEndFlag string,
	nLb []byte,
	fileBlockData []byte) *FileBlock {
	fs := FileBlock{
		Head: FileBlockHead{
			ProtocolType:             rules.FILE_BLOCK_PROTOCOL,
			FileIndexSize:            uint64(len([]byte(fileIndex))), // uuid
			FileNameSize:             uint64(len([]byte(fileName))),
			FileTotalSize:            fileSize,
			FileTotalBlockNum:        fileTotalBlockNum,
			FileBlockPosition:        fileBlockPosition,
			FileBlockSize:            fileBlockSize,
			FileOwnerSize:            fileOwnerSize,
			FileBlockStorageNodeSize: fileBlockStorageNodeSize,
			FileBlockEndFlagSize:     uint32(len([]byte(fileBlockEndFlag))),
		},
		Body: FileBlockBody{
			FileIndex:            []byte(fileIndex),
			FileBlockStorageNode: nLb,
			FileName:             []byte(fileName),
			FileOwner:            []byte(fileOwner),
			FileBlockData:        fileBlockData,
		},
		Foot: FileBlockFoot{
			EndFlag: []byte(fileBlockEndFlag),
		},
	}
	return &fs
}

// create file block arr
func FBNewArrayByFile(file []byte, fileName string, fileOwner string) []*FileBlock {
	fr := bytes.NewReader(file)
	fbuf := make([]byte, len(file))
	fr.Read(fbuf)
	// split file
	fl := utils.SplitFile(fbuf)
	// file position,increment
	FBPosition := uint32(0)
	// file uuid
	fUUID := utils.CrtUUID()
	// file size
	fileSize := uint64(len(file))
	// file total blolck num
	fileTotalBlockNum := fl.Len()
	// ---------------- write processing start ----------------
	fBArr := make([]*FileBlock, fileTotalBlockNum)
	for e := fl.Front(); e != nil; e = e.Next() {
		fileBlock := e.Value.([]byte)
		// create file block storage route table
		fBSNodeRoutable := new(node.TFBRouteTable)
		// encode node list
		nLb, _ := json.Marshal(fBSNodeRoutable.NodeList)
		//// create file block
		fs := FBNew(fUUID, fileName, fileSize, uint64(fileTotalBlockNum), FBPosition, uint32(len(fileBlock)),
			uint32(len([]byte(fileOwner))), uint64(len(nLb)), fileOwner, rules.FILE_BLCOK_END_FLAG, nLb, fileBlock)
		fBArr[FBPosition] = fs
		FBPosition++
	}
	return fBArr
}

// create file block buffer to bytes array
func FBNewByBuf(buff *bytes.Buffer) (*FileBlock, error) {
	fb, err := FBReader(buff)
	if err != nil {
		return nil, err
	}
	return fb, nil
}

// create file block buffer
func FBBuf(fb *FileBlock) (*bytes.Buffer, error) {
	buff := bytes.NewBuffer([]byte{})
	// read in  protocol type
	err := binary.Write(buff, binary.BigEndian, fb.Head.ProtocolType)
	if err != nil {
		return nil, err
	}
	// read in file index
	err = binary.Write(buff, binary.BigEndian, fb.Head.FileIndexSize)
	if err != nil {
		return nil, err
	}
	// read in file name
	err = binary.Write(buff, binary.BigEndian, fb.Head.FileNameSize)
	if err != nil {
		return nil, err
	}
	// read in file total size,unit:byte
	err = binary.Write(buff, binary.BigEndian, fb.Head.FileTotalSize)
	if err != nil {
		return nil, err
	}
	// read in file total block num
	err = binary.Write(buff, binary.BigEndian, fb.Head.FileTotalBlockNum)
	if err != nil {
		return nil, err
	}
	// read in file block position
	err = binary.Write(buff, binary.BigEndian, fb.Head.FileBlockPosition)
	if err != nil {
		return nil, err
	}
	// read in file block size
	err = binary.Write(buff, binary.BigEndian, fb.Head.FileBlockSize)
	if err != nil {
		return nil, err
	}
	// read in file owner size
	err = binary.Write(buff, binary.BigEndian, fb.Head.FileOwnerSize)
	if err != nil {
		return nil, err
	}
	// read in file block storage node size
	err = binary.Write(buff, binary.BigEndian, fb.Head.FileBlockStorageNodeSize)
	if err != nil {
		return nil, err
	}
	// read in file end flag size
	err = binary.Write(buff, binary.BigEndian, fb.Head.FileBlockEndFlagSize)
	if err != nil {
		return nil, err
	}
	// read in file index
	err = binary.Write(buff, binary.BigEndian, fb.Body.FileIndex)
	if err != nil {
		return nil, err
	}
	// read in file block storage node
	err = binary.Write(buff, binary.BigEndian, fb.Body.FileBlockStorageNode)
	if err != nil {
		return nil, err
	}
	// read in file name
	err = binary.Write(buff, binary.BigEndian, fb.Body.FileName)
	if err != nil {
		return nil, err
	}
	// read in file owner
	err = binary.Write(buff, binary.BigEndian, fb.Body.FileOwner)
	if err != nil {
		return nil, err
	}
	// read in file block data
	err = binary.Write(buff, binary.BigEndian, fb.Body.FileBlockData)
	if err != nil {
		return nil, err
	}
	// read in file block end flag
	err = binary.Write(buff, binary.BigEndian, fb.Foot.EndFlag)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

// network unpack by fileblock
func FBNetUnPack(conn net.Conn) (*FileBlock, error) {
	fb := new(FileBlock)
	fb, err := FBReader(conn)
	if err != nil {
		return nil, err
	}
	return fb, err
}

// file block protocol reader
func FBReader(r io.Reader) (*FileBlock, error) {
	fb := new(FileBlock)
	// ---------------------------- protocol head ----------------------------
	// protocol type
	protocolTypeBuf := make([]byte, int(unsafe.Sizeof(fb.Head.ProtocolType)))
	_, err := r.Read(protocolTypeBuf)
	ptBuf := bytes.NewReader(protocolTypeBuf)
	binary.Read(ptBuf, binary.BigEndian, &fb.Head.ProtocolType)
	if err != nil {
		return fb, err
	}
	// file index size
	fileIndexSizeBuf := make([]byte, int(unsafe.Sizeof(fb.Head.FileIndexSize)))
	_, err = r.Read(fileIndexSizeBuf)
	fisBuf := bytes.NewReader(fileIndexSizeBuf)
	binary.Read(fisBuf, binary.BigEndian, &fb.Head.FileIndexSize)
	if err != nil {
		return fb, err
	}
	// file name size
	fileNameSizeBuf := make([]byte, int(unsafe.Sizeof(fb.Head.FileNameSize)))
	_, err = r.Read(fileNameSizeBuf)
	fnsBuf := bytes.NewReader(fileNameSizeBuf)
	binary.Read(fnsBuf, binary.BigEndian, &fb.Head.FileNameSize)
	if err != nil {
		return fb, err
	}
	// file total size
	fileTotalBuf := make([]byte, int(unsafe.Sizeof(fb.Head.FileTotalSize)))
	_, err = r.Read(fileTotalBuf)
	ftBuf := bytes.NewReader(fileTotalBuf)
	binary.Read(ftBuf, binary.BigEndian, &fb.Head.FileTotalSize)
	if err != nil {
		return fb, err
	}
	// file total block num
	fileTotalBlockNumBuf := make([]byte, int(unsafe.Sizeof(fb.Head.FileTotalBlockNum)))
	_, err = r.Read(fileTotalBlockNumBuf)
	ftbBuf := bytes.NewReader(fileTotalBlockNumBuf)
	binary.Read(ftbBuf, binary.BigEndian, &fb.Head.FileTotalBlockNum)
	if err != nil {
		return fb, err
	}
	// file block position
	fileBlockPositionBuf := make([]byte, int(unsafe.Sizeof(fb.Head.FileBlockPosition)))
	_, err = r.Read(fileBlockPositionBuf)
	fbpBuf := bytes.NewReader(fileBlockPositionBuf)
	binary.Read(fbpBuf, binary.BigEndian, &fb.Head.FileBlockPosition)
	if err != nil {
		return fb, err
	}
	// file block size
	fileBlockSizeBuf := make([]byte, int(unsafe.Sizeof(fb.Head.FileBlockSize)))
	_, err = r.Read(fileBlockSizeBuf)
	bsbuf := bytes.NewReader(fileBlockSizeBuf)
	binary.Read(bsbuf, binary.BigEndian, &fb.Head.FileBlockSize)
	if err != nil {
		return fb, err
	}
	// file owner size
	fileOwnerSizeBuf := make([]byte, int(unsafe.Sizeof(fb.Head.FileOwnerSize)))
	_, err = r.Read(fileOwnerSizeBuf)
	fosbuf := bytes.NewReader(fileOwnerSizeBuf)
	binary.Read(fosbuf, binary.BigEndian, &fb.Head.FileOwnerSize)
	if err != nil {
		return fb, err
	}
	// file storage node size
	fileBlockStorageNodeSizeBuf := make([]byte, int(unsafe.Sizeof(fb.Head.FileBlockStorageNodeSize)))
	_, err = r.Read(fileBlockStorageNodeSizeBuf)
	fsnbuf := bytes.NewReader(fileBlockStorageNodeSizeBuf)
	binary.Read(fsnbuf, binary.BigEndian, &fb.Head.FileBlockStorageNodeSize)
	if err != nil {
		return fb, err
	}
	// file end flag size
	fileBlockEndFlagSizeBuf := make([]byte, int(unsafe.Sizeof(fb.Head.FileBlockEndFlagSize)))
	_, err = r.Read(fileBlockEndFlagSizeBuf)
	febuf := bytes.NewReader(fileBlockEndFlagSizeBuf)
	binary.Read(febuf, binary.BigEndian, &fb.Head.FileBlockEndFlagSize)
	if err != nil {
		return fb, err
	}
	// ---------------------------- protocol body ----------------------------
	// file index size
	fileIndexBuf := make([]byte, fb.Head.FileIndexSize)
	n, err := r.Read(fileIndexBuf)
	fb.Body.FileIndex = fileIndexBuf[:n]
	if err != nil {
		return fb, err
	}
	// file block storage node size
	fileStorageNodeBuf := make([]byte, fb.Head.FileBlockStorageNodeSize)
	n, err = r.Read(fileStorageNodeBuf)
	fb.Body.FileBlockStorageNode = fileStorageNodeBuf[:n]
	if err != nil {
		return fb, err
	}
	var s []node.TFBNodeInfo
	json.Unmarshal(fb.Body.FileBlockStorageNode[:], &s)
	// file name data
	fileNameBuf := make([]byte, fb.Head.FileNameSize)
	n, err = r.Read(fileNameBuf)
	fb.Body.FileName = fileNameBuf[:n]
	if err != nil {
		return fb, err
	}
	// file owner data
	fileOwnerBuf := make([]byte, fb.Head.FileOwnerSize)
	n, err = r.Read(fileOwnerBuf)
	fb.Body.FileOwner = fileOwnerBuf[:n]
	if err != nil {
		return fb, err
	}
	// file block data
	fileDataBuf := make([]byte, fb.Head.FileBlockSize)
	n, err = r.Read(fileDataBuf)
	fb.Body.FileBlockData = fileDataBuf[:n]
	if err != nil {
		return fb, err
	}
	// ---------------------------- protocol foot ----------------------------
	// file block end flag data
	fileBlockEndBuf := make([]byte, fb.Head.FileBlockEndFlagSize)
	n, err = r.Read(fileBlockEndBuf)
	fb.Foot.EndFlag = fileBlockEndBuf[:n]
	if err != nil {
		return fb, err
	}
	endFlag := string(fb.Foot.EndFlag)
	if endFlag == rules.FILE_BLCOK_END_FLAG {
		return fb, nil
	}
	return fb, errors.New("illegal agreement")
}

// file block protocol writer
func FBWriter(w io.Writer, fb *FileBlock) error {
	fbArr, err := FBBuf(fb)
	if err != nil {
		return err
	}
	w.Write(fbArr.Bytes())
	return nil
}

// file block encode
func FBEnCoding(fb *FileBlock) ([]byte, error) {
	b, err := json.Marshal(fb)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// file block decode
func FBDnCoding(f []byte) (*FileBlock, error) {
	fb := new(FileBlock)
	err := json.Unmarshal(f, fb)
	if err != nil {
		return nil, err
	}
	return fb, nil
}

// ------------------------- struct handle end -------------------------

// ------------------------- sync handle start -------------------------
func FBSyncFull(fb *FileBlock, netActionType uint8) {
	nodeList := node.TRTNew().NodeList
	for _, n := range nodeList {
		c, err := net.Dial("tcp", n.Addr+":"+n.Port)
		if err != nil {
			// bad node hanlde
		}
		fbb, err := FBBuf(fb)
		np := NPNew(netActionType, fbb.Bytes())
		err = NPWriter(c, np)
		if err != nil {
			// bad node hanlde
		}
	}
}

func FBSyncOne(fb *FileBlock, n node.TNodeInfo, netActionType uint8) {
	c, err := net.Dial("tcp", n.Addr+":"+n.Port)
	if err != nil {
		// bad node hanlde
	}
	fbb, err := FBBuf(fb)
	np := NPNew(netActionType, fbb.Bytes())
	err = NPWriter(c, np)
	if err != nil {
		// bad node hanlde
	}
}

// sync save to localhost
func FBSaveToLocal(fb *FileBlock) {
	fileOwner := string(fb.Body.FileOwner)
	fileName := string(fb.Body.FileName)
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
}

// ------------------------- sync handle end -------------------------

// read file block in local
func RFBByPath(filePath string) (*FileBlock, error) {
	f, err := utils.RFToLocal(filePath)
	if err != nil {
		return nil, err
	}
	fb := new(FileBlock)
	fb, err = FBDnCoding(f)
	if err != nil {
		return nil, err
	}
	return fb, nil
}

// write file to network by fileblock protocol
func WFBToNet(file []byte, fileSize uint64) (bool, error) {
	fr := bytes.NewReader(file)
	fbuf := make([]byte, len(file))
	fr.Read(fbuf)
	// file name
	fileName := test.TFName
	// file account
	ownerAddr := test.TOwner
	// route table
	routeTable := node.TRTNew().NodeList
	// split file
	fl := utils.SplitFile(fbuf)
	// file position,increment
	FBPosition := uint32(0)
	// file uuid
	fUUID := utils.CrtUUID()
	// create channel list
	fBChannel := make(chan FileBlockSyncResult)
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
		go WFBToRt(fUUID, fileName, fileSize, uint64(fl.Len()),
			FBPosition, e.Value.([]byte), ownerAddr, fBSNodeRoutable, fBChannel)
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
func WFBToRt(fUUID string,
	fileName string,
	fileSize uint64,
	fileTotalBlockNum uint64,
	fileBlockPosition uint32,
	fileBlock []byte,
	ownerAddr string,
	fBSNodeRoutable *node.TFBRouteTable,
	c chan FileBlockSyncResult) error {
	// file block result
	fBr := new(FileBlockSyncResult)
	// bad node
	var badNodeList []any
	// encode node list
	nLb, _ := json.Marshal(fBSNodeRoutable.NodeList)
	// create file block
	fs := FBNew(fUUID, fileName, fileSize, fileTotalBlockNum, fileBlockPosition, uint32(len(fileBlock)),
		uint32(len([]byte(ownerAddr))), uint64(len(nLb)), ownerAddr, rules.FILE_BLCOK_END_FLAG, nLb, fileBlock)
	//file block buffer
	buff := bytes.NewBuffer([]byte{})
	// read in file block protocol
	fbbuf, fberr := FBBuf(fs)
	if fberr != nil {
		return fberr
	}
	buff.Write(fbbuf.Bytes())
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
	return nil
}
