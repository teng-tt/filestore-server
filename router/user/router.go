package user

import (
	"filestore-server/api"
	"github.com/gin-gonic/gin"
)

type UserRouter struct {
}

func (u *UserRouter) InitUserRouter(r *gin.Engine) {
	group := r.Group("/user")
	apiGroup := api.ApiGroupApp.UserApiGroup
	group.GET("/signup", apiGroup.UserSignupPage)
	group.POST("/signup", apiGroup.UserSignup)
	group.POST("/signin", apiGroup.UserSignIn)
	group.POST("/info", apiGroup.UserInfo)
}
