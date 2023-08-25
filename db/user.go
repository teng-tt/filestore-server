package db

import (
	"database/sql"
	"filestore-server/global"
	"fmt"
	"log"
)

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

// UserSignup 通过用户名及密码完成user注册
func UserSignup(username, passwd string) bool {
	stmt, err := global.DBConn.Prepare(
		"insert ignore into tbl_user(`user_name`, `user_pwd`) values (?, ?)")
	if err != nil {
		fmt.Println("Failed to insert, er:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert, er:" + err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}
	return false
}

// UserSignIn 判断密码是否一致
func UserSignIn(username, encPwd string) bool {
	stmt, err := global.DBConn.Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println("Failed to select, er:" + err.Error())
		return false
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println("Failed to select, er:" + err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found:" + username)
		return false
	}
	pRows := ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encPwd {
		return true
	}
	return false

}

// UpdateToken 更新用户登录的token
func UpdateToken(username, token string) bool {
	stmt, err := global.DBConn.Prepare(
		"replace into tbl_user_token(`user_name`, `user_token`) values(?,?) ")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func GetUserInfo(username string) (user User, err error) {
	respUser := User{}
	stmt, err := global.DBConn.Prepare(
		"select user_name, signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return respUser, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&respUser.Username, &respUser.SignupAt)
	if err != nil {
		return respUser, err
	}
	return respUser, nil
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
