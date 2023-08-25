package mpupload

import (
	"filestore-server/db"
	"filestore-server/global"
	"filestore-server/model"
	"filestore-server/response"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type MpUploadApi struct {
}

// InitMultipartUpload 初始化分块上传接口
func (m *MpUploadApi) InitMultipartUpload(c *gin.Context) {
	// 1. 解析用户请求参数
	username := c.PostForm("username")
	filehash := c.PostForm("filehash")
	filesize, err := strconv.Atoi(c.PostForm("filesize"))
	if err != nil {
		response.FailWithMessage(c, "params invalid")
		return
	}
	// 2. 获得redis的一个链接
	redisConn := global.CacheConn.Get()
	defer redisConn.Close()

	// 3. 生成分块上传的初始化信息
	uploadInfo := model.MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024, // 5MB
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}
	// 4. 将初始化信息写入到redis缓存
	redisConn.Do("MSET", "MP_"+uploadInfo.UploadID, uploadInfo)
	// 5.将初始化信息返回给客户端
	response.SuccessWithDetailed(c, "ok", uploadInfo)
}

// UploadPart 上传文件分块
func (m *MpUploadApi) UploadPart(c *gin.Context) {
	// 1. 解析用户请求参数
	uploadID := c.PostForm("uploadid")
	chunkIndex := c.PostForm("index")
	// 2. 获得redis连接池中的一个链接
	redisConn := global.CacheConn.Get()
	defer redisConn.Close()
	// 3. 获得文件句柄，用于存储分块内容
	fpath := "/data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		response.FailWithMessage(c, "Upload part failed")
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := c.Request.Body.Read(buf)
		if err != nil {
			break
		}
		fd.Write(buf[:n])
	}
	// 4. 更新redis缓存状态
	redisConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)
	// 5. 返回处理结果到客户端
	response.SuccessWithMessage(c, "OK")
}

// CompleteUpload 通知上传合并接口
func (m *MpUploadApi) CompleteUpload(c *gin.Context) {
	// 1. 解析请求参数
	uploadID := c.PostForm("uploadid")
	username := c.PostForm("username")
	filehash := c.PostForm("filehash")
	filesize, _ := strconv.Atoi(c.PostForm("filesize"))
	filename := c.PostForm("filename")

	// 2. 获得redis连接池中的一个链接
	redisConn := global.CacheConn.Get()
	defer redisConn.Close()

	// 3. 通过uploadid 查询redis并判断是否所有分开上传完成
	data, err := redis.Values(redisConn.Do("HGETALL", "MP+"+uploadID))
	if err != nil {
		response.FailWithMessage(c, "complete upload failed")
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		response.FailWithMessage(c, "invalid request")
		return
	}
	// 4. 合并分开
	// 5. 更新唯一文件表及用户文件表
	db.OnFileUpdateFinished(filehash, filename, "", int64(filesize))
	db.OnUserFileUploadFinished(username, filehash, filename, int64(filesize))

	// 6. 响应处理结果
	response.FailWithMessage(c, "OK")
}

// CanceUploadPart 通知取消上传
func (m *MpUploadApi) CanceUploadPart(c *gin.Context) {
	// 1. 解析用户请求参数：uploadid 上传唯一id
	// 2. 根据uoloadid 删除已经存在的分块文件
	// 3. 根基uploadid 删除redis缓存的状态信息
}

// MultipartUploadStatus 查询上传状态
func (m *MpUploadApi) MultipartUploadStatus(c *gin.Context) {
	// 1. 解析用户请求参数：uploadid 上传唯一id
	// 2. 根据uoloadid 从redis获取分块文件信息
	// 3. 检查分块文件信息的状态是否有效，判断哪些块已经上传
	// 4. 返回以上传完成的块
}
