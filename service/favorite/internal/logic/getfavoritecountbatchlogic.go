package logic

import (
	"context"

	"mini-tiktok/service/favorite/internal/svc"
	"mini-tiktok/service/favorite/pb/favorite"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteCountBatchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteCountBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteCountBatchLogic {
	return &GetFavoriteCountBatchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFavoriteCountBatchLogic) GetFavoriteCountBatch(in *favorite.GetFavoriteCountBatchRequest) (*favorite.GetFavoriteCountBatchResponse, error) {
	batch, missList, err := l.svcCtx.RedisCli.GetFavoriteCountBatch(l.ctx, in.VideoIdList)
	if err != nil {
		return nil, err
	}
	if len(missList) != 0 {
		queryList := make([]int, 0)
		for _, missIdx := range missList {
			queryList = append(queryList, int(in.VideoIdList[missIdx]))
		}
		countBatch, err := l.svcCtx.FavoriteModel.GetVideoFavoriteCountBatch(queryList)
		if err != nil {
			return nil, err
		}
		for i, idx := range missList {
			batch[idx] = countBatch[i]
		}
		//todo 未命中要加入redis中
	}

	res := make([]uint64, 0)
	for _, i := range batch {
		res = append(res, uint64(batch[i]))
	}
	return &favorite.GetFavoriteCountBatchResponse{
		Count: res,
	}, nil
}
