package logic

import (
	"context"

	"mini-tiktok/service/follow/internal/svc"
	"mini-tiktok/service/follow/pb/follow"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendListLogic) GetFriendList(in *follow.GetFriendListRequest) (*follow.GetFriendListResponse, error) {
	resp := new(follow.GetFriendListResponse)
	list, err := l.svcCtx.RelationModel.GetFriendList(uint(in.UserId))
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
