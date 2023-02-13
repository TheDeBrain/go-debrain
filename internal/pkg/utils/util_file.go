package utils

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	"github.com/derain/core/rules"
	"os"
)

// check whether the file exists
func IsFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// check path is exists
func CheckPathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// split file
func SplitFile(fb []byte) list.List {
	// file size
	FILE_SIZE := len(fb)
	// file block size
	MAX_FILE_BLOCK_SIZE := rules.FILE_BLCOK_SIZE_CONSTRAINT
	// file block num, default 1
	var FILE_BLOCK_NUM = int64(1)
	// file size > file block size
	if int64(FILE_SIZE) >= int64(MAX_FILE_BLOCK_SIZE) {
		FILE_BLOCK_NUM = int64(FILE_SIZE) / int64(MAX_FILE_BLOCK_SIZE)
	}
	if (int64(FILE_SIZE) % int64(MAX_FILE_BLOCK_SIZE)) > 0 {
		FILE_BLOCK_NUM += 1
	}
	fileBlockList := list.New()
	// file reader
	r := bytes.NewReader(fb)
	for i := int64(0); i < FILE_BLOCK_NUM; i++ {
		fBuf := make([]byte, MAX_FILE_BLOCK_SIZE)
		fileAt := int64(MAX_FILE_BLOCK_SIZE) * i
		n, err := r.ReadAt(fBuf, fileAt)
		if err != nil {
			fmt.Println(err)
		}
		fBuf = fBuf[:n]
		fileBlockList.PushBack(fBuf)
	}
	return *fileBlockList
}

// read file to localhost
func RFToLocal(filePath string) ([]byte, error) {
	fp, err := os.OpenFile(filePath, os.O_RDONLY, 0777)
	defer fp.Close()
	if err != nil {
		return nil, err
	}
	fInfo, _ := fp.Stat()
	buf := make([]byte, fInfo.Size())
	fp.Read(buf)
	return buf, nil
}

// write file to localhost
func WFToLocal(file []byte, filePath string) {
	f, err3 := os.Create(filePath) //create file
	if err3 != nil {
		fmt.Print(err3)
	}
	_, err := f.Write(file)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	f.Close()
}

// splice file
func SpliceFile(bufSize uint64, fb []byte) *bytes.Buffer {
	sfb := bytes.NewBuffer(make([]byte, bufSize))
	for _, data := range fb {
		binary.Write(sfb, binary.BigEndian, data)
	}
	return sfb
}
