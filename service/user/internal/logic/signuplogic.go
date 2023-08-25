package logic

import (
	"context"
	"filestore-server/db"
	"filestore-server/global"
	"filestore-server/utils"

	"filestore-server/service/user/internal/svc"
	"filestore-server/service/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSignupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignupLogic {
	return &SignupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Signup 处理用户注册请求
func (l *SignupLogic) Signup(in *user.ReqSignup) (*user.RespSignup, error) {
	username := in.Username
	password := in.Password
	resp := new(user.RespSignup)
	if len(username) < 3 || len(password) < 5 {
		resp.Code = -1
		resp.Message = "Invalid parameter!"
		return resp, nil
	}
	encPwd := utils.Sha1([]byte(password + global.PWD_SALT))
	ok := db.UserSignup(username, encPwd)
	if !ok {
		resp.Code = -1
		resp.Message = "注册失败!"
		return resp, nil
	}
	resp.Code = 0
	resp.Message = "success"
	return resp, nil
}
