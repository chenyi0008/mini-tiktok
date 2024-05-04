package logic

import (
	"context"

	"mini-tiktok/service/follow/internal/svc"
	"mini-tiktok/service/follow/pb/follow"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowingListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowingListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowingListLogic {
	return &GetFollowingListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowingListLogic) GetFollowingList(in *follow.GetFollowingListRequest) (*follow.GetFollowingListResponse, error) {
	resp := new(follow.GetFollowingListResponse)
	list, err := l.svcCtx.RelationModel.GetByFollowerId(uint(in.UserId))
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
