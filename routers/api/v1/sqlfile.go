package v1

import (
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	. "github.com/wentaojin/transferdb/conf"
	"github.com/wentaojin/transferdb/server"

	"github.com/wentaojin/transferdb/pkg/e"
	"github.com/wentaojin/transferdb/service"
	"io/ioutil"
	"net/http"
	"strings"
	//"go.uber.org/zap"
	//"github.com/pkg/errors"
)

type sqlfile struct {
	File string `valid:"Required;"`
	service.TargetConfig
}

// 需要知道要执行文件的target
//todo 需要解析一个target
func Sqlfile(c *gin.Context) {
	file := c.Query("file")
	fmt.Println(file)
	valid := validation.Validation{}
	a := sqlfile{File: file}
	ok, _ := valid.Valid(&a)
	//data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if ok {
		// 这里需要 执行文件
		sql, err := ioutil.ReadFile(file)
		if err != nil {
			code = e.ERROR
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": err.Error(),
			})
			return
		}

		err = execsql(string(sql))
		if err != nil {
			code = e.ERROR
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": err.Error(),
			})
			return
		}
		// 执行SQL
		//service.Logger.Info("server run failed",sqls)

		//fmt.Println(cols)
		//fmt.Println(res)

		// 执行sql，统一返回报错

	} else {
		for _, err := range valid.Errors {
			service.Logger.Warn(err.Message)
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
		})
		return
	}
	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
	})
}

func execsql(sqls string) error {
	engine, err := server.NewMySQLEngineGeneralDB(Gcfg.TargetConfig, 300, 10)
	if err != nil {
		return err
	}
	_, _, err = service.Query(engine.MysqlDB, fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s`, conf.Gcfg.TargetConfig.SchemaName))
	if err != nil {
		return err
	}

	sql := strings.Split(string(sqls), ";")

	for _, s := range sql {
		_, _, err = service.Query(engine.MysqlDB, s)
		if err != nil {
			return err
		}
	}
	return nil
}
