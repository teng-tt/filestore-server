package router

import "filestore-server/router/filestore"

type RouterGroup struct {
	FileStoreRouterGroup filestore.FileStoreRouter
}

var RouterGroupApp = new(RouterGroup)
