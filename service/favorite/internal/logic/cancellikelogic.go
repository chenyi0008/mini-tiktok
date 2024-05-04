package logic

import (
	"context"
	"mini-tiktok/common/consts"
	"mini-tiktok/service/favorite/pb/favorite"

	"mini-tiktok/service/favorite/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelLikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelLikeLogic {
	return &CancelLikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelLikeLogic) CancelLike(in *favorite.CancelLikeRequest) (*favorite.Response, error) {
	_, err := l.svcCtx.RedisCli.CancelFavor(l.ctx, in.VideoId, in.UserId)
	//err = l.svcCtx.FavoriteModel.CancelLike(in.UserId, in.VideoId)
	if err != nil {
		return &favorite.Response{
			Code:    consts.FAILED,
			Message: "取消失败",
		}, nil
	}
	return &favorite.Response{
		Code:    consts.SUCCEED,
		Message: "取消成功",
	}, nil
}
