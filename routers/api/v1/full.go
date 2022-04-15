package v1

import (
	"github.com/gin-gonic/gin"
	. "github.com/wentaojin/transferdb/conf"
	"github.com/wentaojin/transferdb/pkg/e"
	"github.com/wentaojin/transferdb/pkg/taskflow"
	"github.com/wentaojin/transferdb/server"
	"github.com/wentaojin/transferdb/service"
	"net/http"
)

//case "full":
//// 全量数据 ETL 非一致性（基于某个时间点，而是直接基于现有 SCN）抽取，离线环境提供与原库一致性
//engine, err := NewEngineDB(
//cfg.SourceConfig, cfg.TargetConfig, cfg.AppConfig.SlowlogThreshold,
//cfg.FullConfig.TableThreads*cfg.FullConfig.SQLThreads*cfg.FullConfig.ApplyThreads)
//if err != nil {
//return err
//}
//if err = taskflow.FullSyncOracleTableRecordToMySQL(cfg, engine); err != nil {
//return err
//}

//包含三块内容， target， source， 以及full的
type Fullform struct {
	FullConfig   service.FullConfig   `form:"full" toml:"full" json:"full"`
	SourceConfig service.SourceConfig `form:"source" toml:"source" json:"source"`
	TargetConfig service.TargetConfig `form:"target" toml:"target" json:"target"`
}

//todo
func Full(c *gin.Context) {
	var form Fullform
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

	// cfg 需要抽离处理
	//todo cfg 需要理解是啥

	Gcfg.TargetConfig = form.TargetConfig
	Gcfg.FullConfig = form.FullConfig
	Gcfg.SourceConfig = form.SourceConfig
	// 	用defer 清空下参数
	//defer func() {
	//	Gcfg.TargetConfig =
	//	Gcfg.FullConfig =
	//	Gcfg.SourceConfig =
	//}()
	//log.Printf("%v",Gcfg)
	engine, err := server.NewEngineDB(
		Gcfg.SourceConfig, Gcfg.TargetConfig, Gcfg.AppConfig.SlowlogThreshold,
		Gcfg.FullConfig.TableThreads*Gcfg.FullConfig.SQLThreads*Gcfg.FullConfig.ApplyThreads)
	if err != nil {
		code = e.ERROR
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": err.Error(),
		})
		return
	}

	if err = taskflow.FullSyncOracleTableRecordToMySQL(Gcfg, engine); err != nil {
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
