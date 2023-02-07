package services

import (
	"github.com/derain/internal/pkg/vars"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSysInfo(c *gin.Context) {
	c.JSON(http.StatusOK, vars.TSys)
}

func GetFileSysInfo(c *gin.Context) {
	c.JSON(http.StatusOK, vars.TFSys)
}


