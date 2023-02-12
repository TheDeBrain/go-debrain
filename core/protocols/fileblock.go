package protocols

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/derain/core/db/table/node"
	"github.com/derain/core/db/table/sys"
	"github.com/derain/core/rules"
	"github.com/derain/internal/pkg/utils"
	"github.com/derain/test"
	"log"
	"net"
	"os"
	"strconv"
)

// Protocol for file block
type FileBlock struct {
	Head FileBlockHead // file block head
	Body FileBlockBody // file block body
	Foot FileBlockFoot // file block foot
}

// Protocol for file block head
type FileBlockHead struct {
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

// create file block pointer
func FBNew(fileIndex string, fileName string, fileSize uint64,
	fileTotalBlockNum uint64, fileBlockPosition uint32,
	fileBlockSize uint32, fileOwnerSize uint32, fileBlockStorageNodeSize uint64,
	ownerAddr string,
	fileBlockEndFlag string,
	nLb []byte,
	fileBlockData []byte) *FileBlock {
	fs := FileBlock{
		Head: FileBlockHead{
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
			FileOwner:            []byte(ownerAddr),
			FileBlockData:        fileBlockData,
		},
		Foot: FileBlockFoot{
			EndFlag: []byte(fileBlockEndFlag),
		},
	}
	return &fs
}

// create file block buffer
func FBBuf(buff *bytes.Buffer, fs *FileBlock) {
	// read in file index
	binary.Write(buff, binary.BigEndian, fs.Head.FileIndexSize)
	// read in file name
	binary.Write(buff, binary.BigEndian, fs.Head.FileNameSize)
	// read in file total size,unit:byte
	binary.Write(buff, binary.BigEndian, fs.Head.FileTotalSize)
	// read in file total block num
	binary.Write(buff, binary.BigEndian, fs.Head.FileTotalBlockNum)
	// read in file block position
	binary.Write(buff, binary.BigEndian, fs.Head.FileBlockPosition)
	// read in file block size
	binary.Write(buff, binary.BigEndian, fs.Head.FileBlockSize)
	// read in file owner size
	binary.Write(buff, binary.BigEndian, fs.Head.FileOwnerSize)
	// read in file block storage node size
	binary.Write(buff, binary.BigEndian, fs.Head.FileBlockStorageNodeSize)
	// read in file end flag size
	binary.Write(buff, binary.BigEndian, fs.Head.FileBlockEndFlagSize)
	// read in file index
	binary.Write(buff, binary.BigEndian, fs.Body.FileIndex)
	// read in file block storage node
	binary.Write(buff, binary.BigEndian, fs.Body.FileBlockStorageNode)
	// read in file name
	binary.Write(buff, binary.BigEndian, fs.Body.FileName)
	// read in file owner
	binary.Write(buff, binary.BigEndian, fs.Body.FileOwner)
	// read in file block data
	binary.Write(buff, binary.BigEndian, fs.Body.FileBlockData)
	// read in file block end flag
	binary.Write(buff, binary.BigEndian, fs.Foot.EndFlag)
}

// network unpack by fileblock
func FBNetUnPack(conn net.Conn) (*FileBlock, error) {
	fb := new(FileBlock)
	// ---------------------------- protocol head ----------------------------
	// file index size
	fileIndexSizeBuf := make([]byte, rules.FILE_INDEX_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err := conn.Read(fileIndexSizeBuf)
	fisBuf := bytes.NewReader(fileIndexSizeBuf)
	binary.Read(fisBuf, binary.BigEndian, &fb.Head.FileIndexSize)
	if err != nil {
		return fb, err
	}
	// file name size
	fileNameSizeBuf := make([]byte, rules.FILE_NAME_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileNameSizeBuf)
	fnsBuf := bytes.NewReader(fileNameSizeBuf)
	binary.Read(fnsBuf, binary.BigEndian, &fb.Head.FileNameSize)
	if err != nil {
		return fb, err
	}
	// file total size
	fileTotalBuf := make([]byte, rules.FILE_TOTAL_SIZE_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileTotalBuf)
	ftBuf := bytes.NewReader(fileTotalBuf)
	binary.Read(ftBuf, binary.BigEndian, &fb.Head.FileTotalSize)
	if err != nil {
		return fb, err
	}
	// file total block num
	fileTotalBlockNumBuf := make([]byte, rules.FILE_TOTAL_BLOCK_NUM_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileTotalBlockNumBuf)
	ftbBuf := bytes.NewReader(fileTotalBlockNumBuf)
	binary.Read(ftbBuf, binary.BigEndian, &fb.Head.FileTotalBlockNum)
	if err != nil {
		return fb, err
	}
	// file block position
	fileBlockPositionBuf := make([]byte, rules.FILE_BLOCK_POSITION_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileBlockPositionBuf)
	fbpBuf := bytes.NewReader(fileBlockPositionBuf)
	binary.Read(fbpBuf, binary.BigEndian, &fb.Head.FileBlockPosition)
	if err != nil {
		return fb, err
	}
	// file block size
	fileBlockSizeBuf := make([]byte, rules.FILE_BLOCK_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileBlockSizeBuf)
	bsbuf := bytes.NewReader(fileBlockSizeBuf)
	binary.Read(bsbuf, binary.BigEndian, &fb.Head.FileBlockSize)
	if err != nil {
		return fb, err
	}
	// file owner size
	fileOwnerSizeBuf := make([]byte, rules.FILE_BLOCK_OWNER_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileOwnerSizeBuf)
	fosbuf := bytes.NewReader(fileOwnerSizeBuf)
	binary.Read(fosbuf, binary.BigEndian, &fb.Head.FileOwnerSize)
	if err != nil {
		return fb, err
	}
	// file storage node size
	fileBlockStorageNodeSizeBuf := make([]byte, rules.FILE_BLOCK_STROAGE_NODE_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileBlockStorageNodeSizeBuf)
	fsnbuf := bytes.NewReader(fileBlockStorageNodeSizeBuf)
	binary.Read(fsnbuf, binary.BigEndian, &fb.Head.FileBlockStorageNodeSize)
	if err != nil {
		return fb, err
	}
	// file end flag size
	fileBlockEndFlagSizeBuf := make([]byte, rules.FILE_BLOCK_END_FLAG_DATASIZE_DESCRIPTOR_BYTE_NUM)
	_, err = conn.Read(fileBlockEndFlagSizeBuf)
	febuf := bytes.NewReader(fileBlockEndFlagSizeBuf)
	binary.Read(febuf, binary.BigEndian, &fb.Head.FileBlockEndFlagSize)
	if err != nil {
		return fb, err
	}
	// ---------------------------- protocol body ----------------------------
	// file index size
	fileIndexBuf := make([]byte, fb.Head.FileIndexSize)
	_, err = conn.Read(fileIndexBuf)
	fb.Body.FileIndex = fileIndexBuf
	if err != nil {
		return fb, err
	}
	fileIndex := string(fileIndexBuf[:])
	fmt.Println(fileIndex)
	// file block storage node size
	fileStorageNodeBuf := make([]byte, fb.Head.FileBlockStorageNodeSize)
	_, err = conn.Read(fileStorageNodeBuf)
	fb.Body.FileBlockStorageNode = fileStorageNodeBuf
	if err != nil {
		return fb, err
	}
	var s []node.TFBNodeInfo
	json.Unmarshal(fb.Body.FileBlockStorageNode[:], &s)
	// file name data
	fileNameBuf := make([]byte, fb.Head.FileNameSize)
	_, err = conn.Read(fileNameBuf)
	fb.Body.FileName = fileNameBuf
	if err != nil {
		return fb, err
	}
	// file owner data
	fileOwnerBuf := make([]byte, fb.Head.FileOwnerSize)
	_, err = conn.Read(fileOwnerBuf)
	fb.Body.FileOwner = fileOwnerBuf
	if err != nil {
		return fb, err
	}
	// file block data
	fileDataBuf := make([]byte, fb.Head.FileBlockSize)
	_, err = conn.Read(fileDataBuf)
	fb.Body.FileBlockData = fileDataBuf
	if err != nil {
		return fb, err
	}
	// ---------------------------- protocol foot ----------------------------
	// file block end flag data
	fileBlockEndBuf := make([]byte, fb.Head.FileBlockEndFlagSize)
	n, err := conn.Read(fileBlockEndBuf)
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

// read file block in local
func RFBInLocal(filePath string) (*FileBlock, error) {
	f, err := utils.RFToLocal(filePath)
	if err != nil {
		return nil, err
	}
	fb := new(FileBlock)
	err = json.Unmarshal(f, fb)
	if err != nil {
		return nil, err
	}
	return fb, nil
}

// read file by FileBlock protocol in net
// Each time a FileBlock is successfully read, a pointer to the FileBlock will be returned
func RFBToLocalInNet(conn net.Conn) error {
	// file block struct
	fb := new(FileBlock)
	fb, err := FBNetUnPack(conn)
	if err != nil {
		return err
	}
	fileOwner := string(fb.Body.FileOwner)
	fileName := string(fb.Body.FileName)
	if string(fb.Foot.EndFlag) == rules.FILE_BLCOK_END_FLAG {
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
func WFBToNet(file []byte, fileSize uint64, ptl *CommProtocol) (bool, error) {
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
func WFBToRt(fUUID string,
	fileName string,
	fileSize uint64,
	fileTotalBlockNum uint64,
	fileBlockPosition uint32,
	fileBlock []byte,
	ownerAddr string,
	fBSNodeRoutable *node.TFBRouteTable,
	ptl *CommProtocol,
	c chan FileBlockSyncResult) {
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
	// read in protocol type
	CPBuf(buff, ptl)
	// read in file block protocol
	FBBuf(buff, fs)
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
