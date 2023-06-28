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

func (u *UserApi) UserSignupPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{})
}

// UserSignup 用户注册
func (u *UserApi) UserSignup(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if len(username) < 3 || len(password) < 5 {
		response.FailWithMessage(c, "Invalid parameter!")
		return
	}
	encPwd := utils.Sha1([]byte(password + PWD_SALT))
	ok := db.UserSignup(username, encPwd)
	if ok {
		response.Success(c)
	} else {
		response.Fail(c)
	}
}

// UserSignIn 用户登录
func (u *UserApi) UserSignIn(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	encPwd := utils.Sha1([]byte(password + PWD_SALT))

	// 1.校验用户名及密码
	pwdChecked := db.UserSignIn(username, encPwd)
	if !pwdChecked {
		response.Fail(c)
		return
	}
	//2. 生成访问凭证token
	token := utils.GetToken(username)
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
	// 拦截器校验
	//token := c.Query("token")
	// 2. 验证token是否有效
	//ok := utils.IsTokenValid(token)
	//if !ok {
	//	c.Status(403)
	//	response.Fail(c)
	//	return
	//}
	// 3. 查询用户信息
	userInfo, err := db.GetUserInfo(username)
	if err != nil {
		c.Status(403)
		response.Fail(c)
	}
	// 4. 返回用户详情
	response.SuccessWithDetailed(c, "OK", userInfo)
}
