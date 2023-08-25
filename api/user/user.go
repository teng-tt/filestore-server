package user

import (
	"filestore-server/db"
	"filestore-server/model"
	"filestore-server/response"
	"filestore-server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	PWD_SALT = "#990"
)

type UserApi struct {
}

// UserSignUpPage 用户注册页面
func (u *UserApi) UserSignUpPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{})
}

// UserSignInPage 用户登录页面
func (u *UserApi) UserSignInPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signin.html", gin.H{})
}

// UserSignUp 用户注册
func (u *UserApi) UserSignUp(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if len(username) < 3 || len(password) < 5 {
		response.FailWithMessage(c, "Invalid parameter!")
		return
	}
	// 将用户密码+密码盐字符串 计算sha1值
	encPwd := utils.Sha1([]byte(password + PWD_SALT))
	// 进行落库
	ok := db.UserSignup(username, encPwd)
	if !ok {
		response.Fail(c)
		return
	}
	response.Success(c)
}

// UserSignIn 用户登录
func (u *UserApi) UserSignIn(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	encPwd := utils.Sha1([]byte(password + PWD_SALT))

	// 1.校验用户名及密码
	pwdChecked := db.UserSignIn(username, encPwd)
	if !pwdChecked {
		response.Fail(c)
		return
	}
	//2. 生成访问凭证token
	token := utils.GetToken(username)
	// 更新用户token记录进行落库
	ok := db.UpdateToken(username, token)
	if !ok {
		response.Fail(c)
		return
	}
	//3. 登录成功重定向到首页
	location := fmt.Sprintf("http://%s/static/view/home.html", c.Request.Host)
	resp := model.UserInfo{
		Location: location,
		Username: username,
		Token:    token,
	}

	response.SuccessWithDetailed(c, "success", resp)
}

func (u *UserApi) UserInfo(c *gin.Context) {
	// 1. 解析请求参数
	username := c.Query("username")
	// 3. 查询用户信息
	userInfo, err := db.GetUserInfo(username)
	if err != nil {
		c.Status(403)
		response.Fail(c)
	}
	// 4. 返回用户详情
	response.SuccessWithDetailed(c, "OK", userInfo)
}
