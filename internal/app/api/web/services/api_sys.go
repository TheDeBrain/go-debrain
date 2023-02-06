package services

import (
	"github.com/derain/core/db/table/sys"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSysInfo(c *gin.Context) {
	sys := new(sys.TSys)
	sys = sys.Load()
	c.JSON(http.StatusOK, sys)
}

func GetFileSysInfo(c *gin.Context) {
	sys := new(sys.TFileSys)
	sys = sys.Load()
	c.JSON(http.StatusOK, sys)
}


