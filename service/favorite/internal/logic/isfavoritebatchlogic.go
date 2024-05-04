package logic

import (
	"context"
	"mini-tiktok/service/favorite/internal/svc"
	"mini-tiktok/service/favorite/pb/favorite"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFavoriteBatchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFavoriteBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFavoriteBatchLogic {
	return &IsFavoriteBatchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFavoriteBatchLogic) IsFavoriteBatch(in *favorite.IsFavoriteBatchRequest) (*favorite.IsFavoriteBatchResponse, error) {
	isFavorite, err := l.svcCtx.RedisCli.GetIsFavoriteBatch(l.ctx, in.IsFavoriteList, in.UserId)
	return &favorite.IsFavoriteBatchResponse{
		IsFavorite: isFavorite,
	}, err

}
