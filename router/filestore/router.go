package filestore

import (
	"filestore-server/api"
	"github.com/gin-gonic/gin"
)

type FileStoreRouter struct {
}

func (f *FileStoreRouter) InitFileStoreRouter(r *gin.Engine) {
	group := r.Group("/file")
	apiGroup := api.ApiGroupApp.FileStoreApiGroup
	group.GET("/upload", apiGroup.FileUploadPage)
	group.POST("/upload", apiGroup.FileUpload)
	group.GET("/meta", apiGroup.GetFileMeta)
	group.GET("/download", apiGroup.FileDownload)
	group.POST("/update", apiGroup.FileMetaUpdate)
	group.DELETE("delete", apiGroup.FileDelete)
}
