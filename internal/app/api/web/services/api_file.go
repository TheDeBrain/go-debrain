package services

import (
	"errors"
	"github.com/derain/core/db/table/sys"
	"github.com/derain/core/rules"
	"github.com/derain/core/sync"
	"github.com/derain/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

// get file
func GetFile(c *gin.Context) error {
	// user address
	fileOwner := c.Query("fileOwner")
	if len(fileOwner) == 0 {
		return errors.New("user owner can not null")
	}
	// file name
	fileName := c.Query("fileName")
	if len(fileName) == 0 {
		return errors.New("user address can not null")
	}
	sync.HandleGetFileBlockReq(fileOwner,fileName)
	return nil
}

// upload file for one
func UploadFileForOne(c *gin.Context) error {
	//fsys := new(sys.TFileSys).Load()
	f, headers, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("Error when try to get file: %v", err)
	}
	//headers.Size
	if headers.Size > rules.MAX_FILE_SIZE {
		log.Println("file size exceeds limit")
		return err
	}
	// file buf
	fbuf := make([]byte, headers.Size)
	f.Read(fbuf)
	bl := utils.SplitFile(fbuf)
	rand.Seed(time.Now().UnixNano())
	for e := bl.Front(); e != nil; e = e.Next() {
		//num := rand.Intn(30001)
		//utils.WFToLocal(e.Value.([]byte), fsys.FileStoragePath+headers.Filename+"-"+string(num))
	}
	c.String(http.StatusOK, headers.Filename)
	// test connect
	conn, _ := net.Dial("tcp", ":"+sys.LoadTSys().SyncPort)
	if conn != nil {
		for e := bl.Front(); e != nil; e = e.Next() {
			err := sync.HandleSendUploadSyncReq(fbuf)
			if err != nil {
			}
		}
	}
	return nil
}

// upload file for more
func UploadFileForMore(c *gin.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
	}
	files := form.File["files"]
	for _, file := range files {
		fSize := file.Size
		fBuf := make([]byte, fSize)
		f, _ := file.Open()
		f.Read(fBuf)
		conn, _ := net.Dial("tcp", ":"+sys.LoadTSys().SyncPort)
		if conn != nil {
			err := sync.HandleSendUploadSyncReq(fBuf)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
