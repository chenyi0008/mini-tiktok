package logic

import (
	"context"
	"mini-tiktok/common/consts"
	"mini-tiktok/service/favorite/pb/favorite"

	"mini-tiktok/service/core/internal/svc"
	"mini-tiktok/service/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeLogic {
	return &LikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeLogic) Like(req *types.LikeRequest) (resp *types.BaseResponse, err error) {
	resp = new(types.BaseResponse)
	var response *favorite.Response
	if req.ActionType == consts.GIVE_LIKE {
		response, err = l.svcCtx.FavoriteRpc.GiveLike(l.ctx, &favorite.GiveLikeRequest{
			UserId:  uint64(req.UserId),
			VideoId: uint64(req.VideoId),
		})
	} else {
		response, err = l.svcCtx.FavoriteRpc.CancelLike(l.ctx, &favorite.CancelLikeRequest{
			UserId:  uint64(req.UserId),
			VideoId: uint64(req.VideoId),
		})
	}

	if err != nil {
		return
	}
	resp.StatusMsg = response.Message
	resp.StatusCode = int(response.Code)
	return
}
