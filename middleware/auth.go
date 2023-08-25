package middleware

import (
	"filestore-server/response"
	"filestore-server/utils"
	"github.com/gin-gonic/gin"
)

// Auth token拦截器验证
func Auth(c *gin.Context) {
	username := c.PostForm("username")
	token := c.PostForm("token")
	// 验证token是否有效
	if len(username) < 3 || !utils.IsTokenValid(token) {
		// 后面的流程不在执行
		c.Abort()
		response.FailWithMessage(c, "token无效")
		return
	}
	// 流转到下一流程
	c.Next()
}
