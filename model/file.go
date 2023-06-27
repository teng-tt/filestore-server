package model

import (
	mydb "filestore-server/db"
)

// FileMeta 文件元信息
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta 新增/更新文件元信息
func UpdateFileMeta(fileMeta FileMeta) {
	fileMetas[fileMeta.FileSha1] = fileMeta
}

// UpdateFileMetaDB 新增/更新文件元信息到mysql中
func UpdateFileMetaDB(fileMeta FileMeta) bool {
	return mydb.OnFileUpdateFinished(fileMeta.FileSha1, fileMeta.FileName, fileMeta.Location, fileMeta.FileSize)
}

// GetFileMeta 通过sha1值获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// GetFIleMetaByDB 从mysql获取文件元信息
func GetFIleMetaByDB(fileHash string) (FileMeta, error) {
	tfile, err := mydb.GetFileMeat(fileHash)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fmeta, nil
}

// RemoveFileMeta 简单删除文件元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
