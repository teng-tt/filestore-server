package api

import "filestore-server/api/filestore"

type ApiGroup struct {
	FileStoreApiGroup filestore.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
