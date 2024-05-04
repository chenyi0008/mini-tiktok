package logic

import (
	"context"
	"mini-tiktok/common/consts"
	"mini-tiktok/service/core/define"
	"mini-tiktok/service/core/helper"

	"mini-tiktok/service/core/internal/svc"
	"mini-tiktok/service/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.UserLoginRequest) (resp *types.UserLoginResponse, err error) {
	resp = new(types.UserLoginResponse)
	user, err := l.svcCtx.UserModel.GetByName(req.Username)
	if err != nil {
		resp.StatusMsg = "不存在此用户"
		resp.StatusCode = consts.FAILED
		return
	}

	if user.Password != helper.Md5(req.Password) {
		resp.StatusMsg = "密码错误"
		resp.StatusCode = consts.FAILED
		return
	}

	resp.Token, err = helper.GenerateToken(user.ID, user.Name, define.TokenExpire)

	resp.UserId = int(user.ID)
	if err != nil {
		resp.StatusCode = consts.FAILED
		resp.StatusMsg = "生成token失败"
		return
	}
	return
}
