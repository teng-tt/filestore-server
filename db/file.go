package db

import (
	"database/sql"
	"filestore-server/global"
	"fmt"
	"log"
)

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// OnFileUpdateFinished 文件上传完成，保存文件元信息
func OnFileUpdateFinished(fileHash, fileName, fileAddr string, filesize int64) bool {
	stmt, err := global.DBConn.Prepare(
		" insert ignore into tbl_file(`file_sha1`, `file_name`, `file_size`," +
			"`file_addr`, `status`) values(?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed tp prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(fileHash, fileName, filesize, fileAddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("file with hash:%s has been uploaded before\n", fileHash)
		}
		return true
	}
	return false
}

// GetFileMeat 从mysql获取文件元信息
func GetFileMeat(fileHash string) (data *TableFile, err error) {
	stmt, err := global.DBConn.Prepare(
		"select file_sha1, file_addr, file_name, file_size from tbl_file" +
			" where file_sha1=? and status=1 limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()
	tfile := TableFile{}
	err = stmt.QueryRow(fileHash).Scan(
		&tfile.FileHash,
		&tfile.FileAddr,
		&tfile.FileName,
		&tfile.FileSize,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &tfile, nil
}

// UpdateFileLocation : 更新文件的存储地址(如文件被转移了)
func UpdateFileLocation(filehash string, fileaddr string) error {
	stmt, err := global.DBConn.Prepare(
		"update tbl_file set`file_addr`=? where  `file_sha1`=? limit 1")
	if err != nil {
		log.Println("预编译sql失败, err:" + err.Error())
		return err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(fileaddr, filehash)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	rf, err := ret.RowsAffected()
	if err != nil {
		return err
	}
	if rf <= 0 {
		log.Printf("更新文件location失败, filehash:%s", filehash)
		err = fmt.Errorf("更新文件location失败, filehash:%s", filehash)
		return err
	}
	return nil
}
