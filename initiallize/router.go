package initiallize

import (
	"filestore-server/middleware"
	"filestore-server/router"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	r := gin.Default()
	// 设置静态资源
	r.Static("/static", "./static")
	// 设置加载模板
	r.LoadHTMLGlob("static/view/**")

	r.Use(middleware.Cors)
	r.Use(middleware.Auth)
	fileStoreRouter := router.RouterGroupApp.FileStoreRouterGroup
	fileStoreRouter.InitFileStoreRouter(r)
	userRouter := router.RouterGroupApp.UserRouterGroup
	userRouter.InitUserRouter(r)

	return r
}
