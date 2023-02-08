package services

import (
	"github.com/derain/core/db/table/sys"
	"github.com/derain/core/sync"
	"github.com/derain/internal/pkg/rules"
	"github.com/derain/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

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
	return nil
}
