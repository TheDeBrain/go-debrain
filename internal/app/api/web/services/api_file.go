package services

import (
	"bytes"
	"errors"
	"github.com/derain/core/db/table/sys"
	"github.com/derain/core/protocols"
	"github.com/derain/core/rules"
	"github.com/derain/core/sync"
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
		c.JSON(http.StatusBadRequest, gin.H{"msg": "user owner can not null"})
		return errors.New("user owner can not null")
	}
	// file name
	fileName := c.Query("fileName")
	if len(fileName) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "user address can not null"})
		return errors.New("user address can not null")
	}
	rc, err := sync.HandleGetFileBlockReq(fileOwner, fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
	}

	a, _ := protocols.RESDecoding(rc.ResultList)

	bb := bytes.NewBuffer(a[0].FileBlock)
	fb, _ := protocols.FBNewByBuf(bb)
	c.JSON(http.StatusOK, gin.H{
		"file_blcok": fb,
	})
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
	// file name
	fileName := c.Request.PostFormValue("fileName")
	// file owner
	fileOwner := c.Request.PostFormValue("fileOwner")
	// file buf
	fbuf := make([]byte, headers.Size)
	n, err := f.Read(fbuf)
	err = sync.HandleSendUploadSyncReq(fbuf[:n], fileName, fileOwner)
	if err != nil {
		return err
	}
	rand.Seed(time.Now().UnixNano())
	c.String(http.StatusOK, headers.Filename)
	return nil
}

// upload file for more
func UploadFileForMore(c *gin.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
	}
	// file array
	files := form.File["files"]
	// file owner
	fileOwner := c.Request.PostFormValue("fileOwner")
	for _, file := range files {
		fileSize := file.Size
		fileName := file.Filename
		fBuf := make([]byte, fileSize)
		f, _ := file.Open()
		f.Read(fBuf)
		conn, _ := net.Dial("tcp", ":"+sys.TSysNew().SyncPort)
		if conn != nil {
			err := sync.HandleSendUploadSyncReq(fBuf, fileName, fileOwner)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
