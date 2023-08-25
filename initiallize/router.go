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
	// 用户管理路由
	userRouter := router.RouterGroupApp.UserRouterGroup
	userRouter.InitUserRouter(r)
	// 文件管理
	fileStoreRouter := router.RouterGroupApp.FileStoreRouterGroup
	fileStoreRouter.InitFileStoreRouter(r)
	// 分块上传
	mpuploadRouter := router.RouterGroupApp.MpUploadRouterGroup
	mpuploadRouter.InitMpUploadRouter(r)

	return r
}
