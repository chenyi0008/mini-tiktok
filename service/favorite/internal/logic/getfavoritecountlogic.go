package logic

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logc"
	"mini-tiktok/service/favorite/favorite"
	"mini-tiktok/service/favorite/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteCountLogic {
	return &GetFavoriteCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFavoriteCountLogic) GetFavoriteCount(in *favorite.GetFavoriteCountRequest) (*favorite.GetFavoriteCountResponse, error) {
	count, err := l.svcCtx.RedisCli.NumOfFavor(l.ctx, in.VideoId)
	if err != nil {
		logc.Error(l.ctx, err.Error())
	}
	if err == redis.Nil {
		count, err2 := l.svcCtx.FavoriteModel.CountByVideoId(uint(in.VideoId))
		logc.Info(l.ctx, "2.count: ", count, "  err: ", err)
		return &favorite.GetFavoriteCountResponse{
			Count: uint64(count),
		}, err2
	}

	return &favorite.GetFavoriteCountResponse{
		Count: uint64(count),
	}, err

}
