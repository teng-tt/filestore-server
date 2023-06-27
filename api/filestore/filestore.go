package filestore

import (
	"encoding/json"
	filemeta "filestore-server/model"
	"filestore-server/response"
	"filestore-server/utils"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type FileStoreApi struct {
}

func (f *FileStoreApi) FileUploadPage(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", gin.H{})
}

// FileUpload 文件上传
func (f *FileStoreApi) FileUpload(c *gin.Context) {
	// 接收文件流存储在本地
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage(c, "Failed to get data, err："+err.Error())
		return
	}
	src, err := file.Open()
	if err != nil {
		response.FailWithMessage(c, "upload file failed, err："+err.Error())
		return
	}
	defer src.Close()
	fileMeta := filemeta.FileMeta{
		FileName: file.Filename,
		Location: "/tmp/" + file.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		response.FailWithMessage(c, "upload file failed, err："+err.Error())
		return
	}

	defer newFile.Close()
	fileMeta.FileSize, err = io.Copy(newFile, src)
	if err != nil {
		response.FailWithMessage(c, "upload file failed, err："+err.Error())
		return
	}
	newFile.Seek(0, 0)
	fileMeta.FileSha1 = utils.FileSha1(newFile)
	// 更新文件元信息，写入内存变量
	// filemeta.UpdateFileMeta(fileMeta)
	filemeta.UpdateFileMetaDB(fileMeta)
	response.Success(c)
}

// GetFileMeta 获取文件元信息
func (f *FileStoreApi) GetFileMeta(c *gin.Context) {
	fileHash := c.Query("filehash")
	fileMeta, err := filemeta.GetFIleMetaByDB(fileHash)
	if err != nil {
		response.Fail(c)
		return
	}
	data, err := json.Marshal(fileMeta)
	if err != nil {
		response.Fail(c)
		return
	}
	response.SuccessWithDetailed(c, "success!", data)
}

// FileDownload 文件下载
func (f *FileStoreApi) FileDownload(c *gin.Context) {
	fileHash := c.Query("filehash")
	fm := filemeta.GetFileMeta(fileHash)
	file, err := os.Open(fm.Location)
	if err != nil {
		response.Fail(c)
		return
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		response.Fail(c)
		return
	}
	c.Header("Content-Type", "application/octect-stream")
	c.Header("content-disposition", "attachment;filename="+fm.FileName)
	response.SuccessWithDetailed(c, "success!", data)
}

// FileMetaUpdate 更新文件元信息
func (f *FileStoreApi) FileMetaUpdate(c *gin.Context) {
	optYpe := c.Query("op")
	fileSha1 := c.Query("filehash")
	newFileName := c.Query("filename")

	if optYpe != "0" {
		c.Status(403)
		response.Fail(c)
		return
	}
	if optYpe != "POST" {
		c.Status(405)
		response.Fail(c)
		return
	}

	currFileMeta := filemeta.GetFileMeta(fileSha1)
	currFileMeta.FileName = newFileName
	filemeta.UpdateFileMeta(currFileMeta)
	data, err := json.Marshal(currFileMeta)
	if err != nil {
		c.Status(500)
		response.Fail(c)
	}
	response.SuccessWithDetailed(c, "success!", data)
}

// FileDelete 删除文件
func (f *FileStoreApi) FileDelete(c *gin.Context) {
	fileSha1 := c.Query("filehash")
	fm := filemeta.GetFileMeta(fileSha1)
	// 删除文件
	os.Remove(fm.Location)
	// 删除文件元信息
	filemeta.RemoveFileMeta(fm.FileSha1)
	response.Success(c)
}
