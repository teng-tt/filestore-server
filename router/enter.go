package router

import (
	"filestore-server/router/filestore"
	"filestore-server/router/user"
)

type RouterGroup struct {
	FileStoreRouterGroup filestore.FileStoreRouter
	UserRouterGroup      user.UserRouter
}

var RouterGroupApp = new(RouterGroup)
