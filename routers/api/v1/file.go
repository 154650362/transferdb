package v1

import "github.com/gin-gonic/gin"

func File(c *gin.Context) {
	path := ""
	fileName := path + c.Query("name")
	c.File(fileName)
}
