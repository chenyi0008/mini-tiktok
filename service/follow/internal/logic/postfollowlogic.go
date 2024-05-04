package logic

import (
	"context"

	"mini-tiktok/service/follow/internal/svc"
	"mini-tiktok/service/follow/pb/follow"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostFollowLogic {
	return &PostFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PostFollowLogic) PostFollow(in *follow.PostFollowRequest) (*follow.Response, error) {
	resp := new(follow.Response)
	var err error

	if in.ActionType == 1 {
		_, err = l.svcCtx.RelationModel.Create(uint(in.UserId), uint(in.ToUserId), in.Username)
	} else {
		err = l.svcCtx.RelationModel.DeleteById(uint(in.UserId), uint(in.ToUserId))
	}

	if err != nil {
		resp.Message = "请求失败"
		resp.Code = 1
		return resp, nil
	}
	resp.Message = "请求成功"
	return resp, nil
}
