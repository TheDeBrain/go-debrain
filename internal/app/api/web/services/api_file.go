package services

import (
	"bytes"
	"errors"
	"github.com/derain/core/db/table/sys"
	"github.com/derain/core/protocols"
	"github.com/derain/core/rules"
	"github.com/derain/core/sync/file"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

// get file
func GetFile(c *gin.Context, netType string) error {
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
	rc, err := file.HandleGetFileBlockReqTCP(fileOwner, fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
	}

	a, _ := protocols.RESDecoding(rc.ResultList)

	fbs, err := protocols.FBReaderMore(bytes.NewReader(a[0].FileBlock))

	//bb := bytes.NewBuffer(a[0].FileBlock)
	//fb, _ := protocols.FBNewByBuf(bb)
	c.JSON(http.StatusOK, gin.H{
		"file_blcok": fbs,
	})
	return nil
}

// upload file for one
func UploadFileForOne(c *gin.Context, netType string) error {
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
	err = file.HandleSendUploadSyncReqTCP(fbuf[:n], fileName, fileOwner)
	if err != nil {
		return err
	}
	rand.Seed(time.Now().UnixNano())
	c.String(http.StatusOK, headers.Filename)
	return nil
}

// upload file for more
func UploadFileForMore(c *gin.Context, netType string) error {
	form, err := c.MultipartForm()
	if err != nil {
	}
	// file array
	files := form.File["files"]
	// file owner
	fileOwner := c.Request.PostFormValue("fileOwner")
	for _, fl := range files {
		fileSize := fl.Size
		fileName := fl.Filename
		fBuf := make([]byte, fileSize)
		f, _ := fl.Open()
		f.Read(fBuf)
		conn, _ := net.Dial("udp", ":"+string(sys.TSysNew().SyncPortUDP))
		if conn != nil {
			err := file.HandleSendUploadSyncReqTCP(fBuf, fileName, fileOwner)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
