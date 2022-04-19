package v1

import (
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	//"github.com/wentaojin/transferdb/conf"
	"github.com/wentaojin/transferdb/pkg/e"
	//"github.com/wentaojin/transferdb/server"
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
func Sqlfile(c *gin.Context) {
	file := c.Query("file")
	fmt.Println(file)
	valid := validation.Validation{}
	a := sqlfile{File: file}
	ok, _ := valid.Valid(&a)
	data := make(map[string]interface{})
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

		//engine, err := server.NewMySQLEngineGeneralDB(conf.Gcfg.TargetConfig, 300, 10)
		//if err != nil {
		//	service.Logger.Info("server run failed", zap.Error(errors.Cause(err)))
		//}
		sqls := strings.Split(string(sql), ";")
		for _, sql := range sqls {
			fmt.Println("------------------------------")
			fmt.Println(sql)
			//fmt.Println("------------------------------")
			//cols, res, err := service.Query(engine.MysqlDB, "select * from t1")
			//if err != nil {
			//	fmt.Println(err)
			//}
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
		// 就要参数失败
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
		})
		return
	}
	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
