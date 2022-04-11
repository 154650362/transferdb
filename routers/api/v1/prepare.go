package v1

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

import (
	"github.com/gin-gonic/gin"
	. "github.com/wentaojin/transferdb/conf"
	"github.com/wentaojin/transferdb/pkg/e"
	"github.com/wentaojin/transferdb/pkg/prepare"
	"github.com/wentaojin/transferdb/server"
	"github.com/wentaojin/transferdb/service"
	"net/http"
)

//todo ? 是否可以使用统一的一个config的结构
// 需要替换掉， 先做那个必要配置文件的初始化
// 这是一个demo， 需要修改
//type PrepareForm struct {
//	Username   string `form:"username" json:"username" uri:"username" xml:"username" binding:"required"`
//	Password   string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
//	RePassword string `form:"rePassword" json:"rePassword" uri:"rePassword" xml:"rePassword" binding:"required"`
//	Nickname   string `form:"nickname" json:"nickname" uri:"nickname" xml:"nickname" binding:"required"`
//	Captcha    string `form:"captcha" json:"captcha" uri:"captcha" xml:"captcha" binding:"required"`
//}

//type TargetConfig struct {
//	Username      string `form:"username" toml:"username" json:"username" binding:"required"`
//	Password      string `form:"password" toml:"password" json:"password" binding:"required"`
//	Host          string `form:"host" toml:"host" json:"host" binding:"required"`
//	Port          int    `form:"port" toml:"port" json:"port" binding:"required"`
//	ConnectParams string `form:"connect-params" toml:"connect-params" toml:"connect-params" json:"connect-params" binding:"required"`
//	MetaSchema    string `form:"meta-schema" toml:"meta-schema" json:"meta-schema"`
//	SchemaName    string `form:"schema-name" toml:"schema-name",json:"schema-name" binding:"required"`
//	Overwrite     bool   `form:"overwrite" toml:"overwrite" json:"overwrite"`
//}

func Prepare(c *gin.Context) {
	var code int
	code = e.SUCCESS

	var form service.TargetConfig
	if err := c.Bind(&form); err != nil {
		code = e.INVALID_PARAMS
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": err.Error(),
		})
		return
	}

	//数据格式正确
	// todo进行 调用真正的执行函数
	//case "prepare":
	//	// 表结构转换 - only prepare 阶段

	// cfg 需要抽离处理
	//todo cfg 需要理解是啥
	//var cfg *service.CfgFile

	Gcfg.TargetConfig = form
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
