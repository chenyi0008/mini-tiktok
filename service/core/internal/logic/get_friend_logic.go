package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"mini-tiktok/service/follow/pb/follow"

	"mini-tiktok/service/core/internal/svc"
	"mini-tiktok/service/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendLogic {
	return &GetFriendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendLogic) GetFriend(req *types.GetFriendRequest) (resp *types.GetFriendResponse, err error) {
	// todo: add your logic here and delete this line
	resp = new(types.GetFriendResponse)
	result, err := l.svcCtx.FollowRpc.GetFriendList(l.ctx, &follow.GetFriendListRequest{
		UserId: uint64(req.UserId),
	})
	logx.Error(err)
	if err != nil {
		resp.StatusCode = 1
	}
	copier.Copy(&resp.UserList, result.UserList)

	return
}
