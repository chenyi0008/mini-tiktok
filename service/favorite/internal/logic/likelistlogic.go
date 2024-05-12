package logic

import (
	"context"
	"mini-tiktok/common/consts"
	"mini-tiktok/service/favorite/internal/svc"
	"mini-tiktok/service/favorite/pb/favorite"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeListLogic {
	return &LikeListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LikeListLogic) LikeList(in *favorite.LikeListRequest) (*favorite.LikeListResponse, error) {
	resp := new(favorite.LikeListResponse)
	resList, err := l.svcCtx.RedisCli.GetFavorSet(l.ctx, in.UserId)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	if len(resList) == 0 {
		list, err := l.svcCtx.FavoriteModel.GetByUserId(in.UserId)
		logx.Info(list)
		if err != nil {
			logx.Error(err)
			return nil, err
		}
		for _, item := range *list {
			resList = append(resList, uint64(item.VideoId))
		}
		err = l.svcCtx.RedisCli.AddFavorSet(l.ctx, in.UserId, resList)
		if err != nil {
			logx.Error(err)
			return nil, err
		}
	}

	return &favorite.LikeListResponse{
		Code:    consts.SUCCEED,
		VideoId: resList,
		Message: "查询成功",
	}, nil
	return resp, nil
}
