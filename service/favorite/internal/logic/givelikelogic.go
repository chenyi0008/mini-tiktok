package logic

import (
	"context"
	"mini-tiktok/common/consts"
	"mini-tiktok/service/favorite/pb/favorite"

	"mini-tiktok/service/favorite/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GiveLikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGiveLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GiveLikeLogic {
	return &GiveLikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GiveLikeLogic) GiveLike(in *favorite.GiveLikeRequest) (*favorite.Response, error) {
	_, err := l.svcCtx.RedisCli.SetFavor(l.ctx, in.VideoId, in.UserId)
	if err != nil {
		logx.Error(err)
		return &favorite.Response{
			Code:    consts.FAILED,
			Message: "点赞失败",
		}, nil
	}

	//err = l.svcCtx.FavoriteModel.GiveLike(in.UserId, in.VideoId)
	if err != nil {
		logx.Error(err)
		return &favorite.Response{
			Code:    consts.FAILED,
			Message: "点赞失败",
		}, nil
	}
	return &favorite.Response{
		Code:    0,
		Message: "点赞成功",
	}, nil
}
