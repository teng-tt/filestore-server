package filestore

import (
	"encoding/json"
	"filestore-server/db"
	"filestore-server/global"
	filemeta "filestore-server/model"
	"filestore-server/mq"
	"filestore-server/response"
	"filestore-server/store/oss"
	"filestore-server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
		Location: "tmp/" + file.Filename,
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

	newFile.Seek(0, 0)
	// 同时将文件写入ceph存储
	//data, _ := ioutil.ReadAll(newFile)
	//bucket := ceph.GetCephBucket("userfile")
	//cephPath := "/ceph/" + fileMeta.FileSha1
	//_ = bucket.Put(cephPath, data, "octet-stream", s3.PublicRead)
	//fileMeta.Location = cephPath

	// 同时将文件写入到oss
	ossPath := "oss/" + fileMeta.FileSha1
	//err = oss.Bucket().PutObject(ossPath, newFile)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	response.FailWithMessage(c, "Upload failed!")
	//	return
	//}
	//fileMeta.Location = ossPath

	// 加入消息队列异步转存到oss
	data := mq.TransferData{
		FileHash:     fileMeta.FileSha1,
		CurLocation:  fileMeta.Location,
		DestLocation: ossPath,
	}
	pubData, _ := json.Marshal(data)
	ok := mq.Publish(
		global.CONF.RabbitConf.TransExchangeName,
		global.CONF.RabbitConf.TransOSSRoutingKey,
		pubData)
	if !ok {
		// todo 加入重试发送消息逻辑
	}
	// 更新文件元信息，写入内存变量
	// filemeta.UpdateFileMeta(fileMeta)
	filemeta.UpdateFileMetaDB(fileMeta)
	// 更新用户文件表记录
	username := c.Query("username")
	if ok := db.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize); ok {
		response.Success(c)
	} else {
		response.FailWithMessage(c, "Upload Failed")
	}
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

// GetUserFileMetas 批量获取用户文件元信息
func (f *FileStoreApi) GetUserFileMetas(c *gin.Context) {
	username := c.Query("username")
	limit, _ := strconv.Atoi(c.PostForm("limit"))
	userFiles, err := db.QueryUserFileMetas(username, limit)
	if err != nil {
		response.Fail(c)
		return
	}
	response.SuccessWithDetailed(c, "OK", userFiles)
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

// TryFastUpload 尝试秒传接口
func (f *FileStoreApi) TryFastUpload(c *gin.Context) {
	// 1. 解析请求参数
	username := c.PostForm("username")
	filehash := c.PostForm("filehash")
	filename := c.PostForm("filename")
	filesize, _ := strconv.ParseInt(c.PostForm("filesize"), 10, 64)

	// 2. 从文件表中查询相同的hash的文件记录
	fileMeta, err := db.GetFileMeat(filehash)
	if err != nil {
		fmt.Println(err.Error())
		response.Fail(c)
		return
	}

	// 3. 查不到记录则返回秒传失败
	if fileMeta == nil {
		response.FailWithMessage(c, "秒传失败，请访问葡萄上传接口")
		return
	}
	// 4. 上传过则将文件信息吸入用户文件表，返回成功
	ok := db.OnUserFileUploadFinished(username, filehash, filename, filesize)
	if ok {
		response.SuccessWithMessage(c, "秒传成功！")
		return
	}
	response.FailWithMessage(c, "秒传失败，请稍后重试！")
}

// DownloadURL 生成oss文件的下载地址
func (f *FileStoreApi) DownloadURL(c *gin.Context) {
	filehash := c.PostForm("filehash")
	// 从文件表查找记录
	row, _ := db.GetFileMeat(filehash)
	// todo 判断文件是存在 OSS 还是Ceph
	signedURL := oss.DownloadURL(row.FileAddr.String)
	response.SuccessWithDetailed(c, "OK", signedURL)
}
