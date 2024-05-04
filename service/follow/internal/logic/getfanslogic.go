package logic

import (
	"context"

	"mini-tiktok/service/follow/internal/svc"
	"mini-tiktok/service/follow/pb/follow"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFansLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFansLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFansLogic {
	return &GetFansLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFansLogic) GetFans(in *follow.GetFansRequest) (*follow.GetFansResponse, error) {
	resp := new(follow.GetFansResponse)
	list, err := l.svcCtx.RelationModel.GetByFollowingId(uint(in.UserId))
	if err != nil {
		resp.Code = 1
		return resp, nil
	}
	for _, relation := range list {
		user := &follow.User{
			Id:   uint64(relation.FollowingId),
			Name: relation.FollowingName,
		}
		resp.UserList = append(resp.UserList, user)
	}
	return resp, nil
}
