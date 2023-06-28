package middleware

import (
	"filestore-server/response"
	"filestore-server/utils"
	"github.com/gin-gonic/gin"
)

// Auth token拦截器验证
func Auth(c *gin.Context) {
	username := c.Query("username")
	token := c.Query("token")
	if len(username) < 3 || !utils.IsTokenValid(token) {
		response.Fail(c)
		return
	}
	c.Next()
}
