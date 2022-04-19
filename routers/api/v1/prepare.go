package v1

import (
	"github.com/gin-gonic/gin"
	. "github.com/wentaojin/transferdb/conf"
	"github.com/wentaojin/transferdb/pkg/e"
	"github.com/wentaojin/transferdb/pkg/prepare"
	"github.com/wentaojin/transferdb/server"
	"github.com/wentaojin/transferdb/service"
	"log"
	"net/http"
)

func Prepare(c *gin.Context) {
	var form service.TargetConfig
	var code int
	code = e.SUCCESS

	if err := c.Bind(&form); err != nil {
		code = e.INVALID_PARAMS
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": err.Error(),
		})
		return
	}

	//var cfg *service.CfgFile

	Gcfg.TargetConfig = form
	// debug使用， 后面删掉
	log.Printf("%v", Gcfg)

	engine, err := server.NewMySQLEnginePrepareDB(Gcfg.TargetConfig, Gcfg.AppConfig.SlowlogThreshold, 1024)
	if err != nil {
		code = e.ERROR
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": err.Error(),
		})
		return
	}

	if err = prepare.TransferDBEnvPrepare(engine); err != nil {
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
