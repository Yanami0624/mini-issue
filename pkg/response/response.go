package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, data any) {
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

func Fail(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{
		"code": -1,
		"msg":  msg,
		"data": nil,
	})
}
