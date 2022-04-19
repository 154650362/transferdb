package v1

import (
	"github.com/gin-gonic/gin"
	"os"
)

func File(c *gin.Context) {
	path, _ := os.Getwd()
	fileName := path + c.Query("name")
	c.File(fileName)
}

//func DownloadFileService(c *gin.Context) {
//	fileDir := c.Query("fileDir")
//	fileName := c.Query("fileName")
//	//打开文件
//	_, errByOpenFile := os.Open(fileDir + "/" + fileName)
//	//非空处理
//	if common.IsEmpty(fileDir) || common.IsEmpty(fileName) || errByOpenFile != nil {
//		/*c.JSON(http.StatusOK, gin.H{
//		    "success": false,
//		    "message": "失败",
//		    "error":   "资源不存在",
//		})*/
//		c.Redirect(http.StatusFound, "/404")
//		return
//	}
//	c.Header("Content-Type", "application/octet-stream")
//	c.Header("Content-Disposition", "attachment; filename="+fileName)
//	c.Header("Content-Transfer-Encoding", "binary")
//	c.File(fileDir + "/" + fileName)
//	return
//}
