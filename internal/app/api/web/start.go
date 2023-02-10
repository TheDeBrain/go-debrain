package web

import (
	"github.com/derain/core/db/table/sys"
	"github.com/derain/internal/app/api/web/services"
	"github.com/gin-gonic/gin"
)

func StartWebApiService() {
	r := gin.Default()
	// init server
	initServer(r)
	// set port
	r.Run(":" + sys.LoadTSys().WebApiPort)
}

func initServer(r *gin.Engine) {
	// get system info
	r.GET("/sys/getSysInfo", func(c *gin.Context) {
		services.GetSysInfo(c)
	})
	// get file system info
	r.GET("/sys/getFileSysInfo", func(c *gin.Context) {
		services.GetFileSysInfo(c)
	})
	// get file
	r.GET("/file/getFile", func(c *gin.Context) {
		services.GetFile(c)
	})
	// upload file for one
	r.POST("/file/upload/one", func(c *gin.Context) {
		services.UploadFileForOne(c)
	})
	// upload file for more
	r.POST("/file/upload/more", func(c *gin.Context) {
		services.UploadFileForMore(c)
	})
}
