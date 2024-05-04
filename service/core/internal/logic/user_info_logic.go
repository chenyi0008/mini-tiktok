package logic

import (
	"context"
	"mini-tiktok/service/core/helper"

	"mini-tiktok/service/core/internal/svc"
	"mini-tiktok/service/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	resp = new(types.UserInfoResponse)
	user, err := helper.AnalyzeToken(req.Token)
	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = "token无效"
		return
	}
	resp.User.Name = user.Username
	resp.User.ID = int(user.Id)

	return
}
