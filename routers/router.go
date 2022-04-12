package routers

import (
	"github.com/gin-gonic/gin"
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

	apiv1 := r.Group("/api/v1")
	{

		//health 状态
		apiv1.GET("/test", v1.Test)

		apiv1.POST("/prepare", v1.Prepare)

		apiv1.POST("/full", v1.Full)
		//

	}
	return r
}
