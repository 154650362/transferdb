package v1

import (
	"github.com/gin-gonic/gin"
	. "github.com/wentaojin/transferdb/conf"
	"github.com/wentaojin/transferdb/pkg/e"
	"github.com/wentaojin/transferdb/pkg/reverser"
	"github.com/wentaojin/transferdb/server"
	"github.com/wentaojin/transferdb/service"
	"net/http"
)

//包含2块内容， target， source，
type Reverseform struct {
	SourceConfig service.SourceConfig `form:"source" toml:"source" json:"source"`
	TargetConfig service.TargetConfig `form:"target" toml:"target" json:"target"`
}

func Reverse(c *gin.Context) {
	var form Reverseform
	var code int
	code = e.SUCCESS

	if err := c.BindJSON(&form); err != nil {
		code = e.INVALID_PARAMS
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": err.Error(),
		})
		return
	}

	Gcfg.TargetConfig = form.TargetConfig
	Gcfg.SourceConfig = form.SourceConfig

	engine, err := server.NewEngineDB(Gcfg.SourceConfig, Gcfg.TargetConfig,
		Gcfg.AppConfig.SlowlogThreshold, 1024)
	if err != nil {
		code = e.ERROR
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": err.Error(),
		})
		return
	}

	if err = reverser.ReverseOracleToMySQLTable(engine, Gcfg); err != nil {
		code = e.ERROR
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
	})
}
