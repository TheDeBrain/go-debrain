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
	r.Run(":" + string(sys.TSysNew().WebApiPort))
}

func initServer(r *gin.Engine) {
	// ------------------ system api start ------------------
	// get system info
	r.GET("/sys/getSysInfo", func(c *gin.Context) {
		services.GetSysInfo(c)
	})
	// get file system info
	r.GET("/sys/getFileSysInfo", func(c *gin.Context) {
		services.GetFileSysInfo(c)
	})
	// ------------------ system api end ------------------

	// ------------------ file api start ------------------
	// get file
	r.GET("/file/getFile", func(c *gin.Context) {
		services.GetFile(c,"tcp")
	})
	// upload file for one
	r.POST("/file/upload/one", func(c *gin.Context) {
		services.UploadFileForOne(c,"udp")
	})
	// upload file for more
	r.POST("/file/upload/more", func(c *gin.Context) {
		services.UploadFileForMore(c,"udp")
	})
	// ------------------ file api end ------------------
}
