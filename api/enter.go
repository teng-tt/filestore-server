package api

import (
	"filestore-server/api/filestore"
	"filestore-server/api/mpupload"
	"filestore-server/api/user"
)

type ApiGroup struct {
	FileStoreApiGroup filestore.ApiGroup
	UserApiGroup      user.ApiGroup
	MpUploadApiGroup  mpupload.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
