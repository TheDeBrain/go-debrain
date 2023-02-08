package services

import (
	"github.com/derain/core/db/table/sys"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSysInfo(c *gin.Context) {
	c.JSON(http.StatusOK, sys.LoadTSys())
}

func GetFileSysInfo(c *gin.Context) {
	c.JSON(http.StatusOK, sys.LoadFileSys())
}


