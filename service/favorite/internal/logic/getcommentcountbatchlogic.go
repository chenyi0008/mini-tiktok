package logic

import (
	"context"

	"mini-tiktok/service/favorite/internal/svc"
	"mini-tiktok/service/favorite/pb/favorite"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentCountBatchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentCountBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentCountBatchLogic {
	return &GetCommentCountBatchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentCountBatchLogic) GetCommentCountBatch(in *favorite.GetCommentCountBatchRequest) (*favorite.GetCommentCountBatchResponse, error) {
	batch, missList, err := l.svcCtx.RedisCli.GetCommentCountBatch(l.ctx, in.VideoIdList)
	if err != nil {
		return nil, err
	}
	if len(missList) != 0 {
		queryList := make([]int, 0)
		for _, missIdx := range missList {
			queryList = append(queryList, int(in.VideoIdList[missIdx]))
		}
		// 从mysql读数据
		countBatch, err := l.svcCtx.CommentModel.GetVideoCommentCountBatch(queryList)
		if err != nil {
			return nil, err
		}
		for i, idx := range missList {
			batch[idx] = countBatch[i]
			// 把数据存到redis
			_, err := l.svcCtx.RedisCli.SetVideoCommentCount(l.ctx, queryList[i], countBatch[i])
			if err != nil {
				return nil, err
			}
		}
	}

	res := make([]uint64, 0)
	for _, i := range batch {
		res = append(res, uint64(batch[i]))
	}

	return &favorite.GetCommentCountBatchResponse{
		Count: res,
	}, nil
}
