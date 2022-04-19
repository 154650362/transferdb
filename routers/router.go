package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/wentaojin/transferdb/middleware/jwt"
	v1 "github.com/wentaojin/transferdb/routers/api/v1"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode("debug")
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.GET("/auth", GetAuth)

	apiv1 := r.Group("/api/v1").Use(jwt.JWT())
	{

		apiv1.GET("/test", v1.Test)
		//做其他任务必须先做prepare
		//todo 待测试
		apiv1.POST("/prepare", v1.Prepare)
		//todo

		//todo 待测试
		apiv1.POST("/full", v1.Full)

		//todo 待测试
		apiv1.POST("/reverse", v1.Reverse)
		//file server,使用query参数传递下载文件的参数
		//sql 执行接口，创建数据库执行sql文件
		apiv1.GET("/sqlfile", v1.Sqlfile)
		apiv1.GET("/file", v1.File)
	}
	return r
}
