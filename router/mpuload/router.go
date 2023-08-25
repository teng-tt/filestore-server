package mpuload

import "github.com/gin-gonic/gin"
import "filestore-server/api"

type MpUploadRouter struct {
}

func (m *MpUploadRouter) InitMpUploadRouter(r *gin.Engine) {
	// 分块上传接口
	group := r.Group("/file/mpupload")
	apiGroup := api.ApiGroupApp.MpUploadApiGroup
	group.POST("/init", apiGroup.InitMultipartUpload)
	group.POST("/uppart", apiGroup.UploadPart)
	group.POST("/complete", apiGroup.CompleteUpload)
}
