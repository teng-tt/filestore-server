package router

import (
	"filestore-server/router/filestore"
	"filestore-server/router/mpuload"
	"filestore-server/router/user"
)

type RouterGroup struct {
	FileStoreRouterGroup filestore.FileStoreRouter
	UserRouterGroup      user.UserRouter
	MpUploadRouterGroup  mpuload.MpUploadRouter
}

var RouterGroupApp = new(RouterGroup)
