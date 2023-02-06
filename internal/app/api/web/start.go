package web

import (
	"github.com/derain/core/db/table/sys"
	"github.com/derain/internal/app/api/web/services"
	"github.com/gin-gonic/gin"
)

func StartWebApiServer() {
	r := gin.Default()
	// init server
	initServer(r)
	// set port
	r.Run(":" + new(sys.TSys).Load().WebApiPort)
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
}
