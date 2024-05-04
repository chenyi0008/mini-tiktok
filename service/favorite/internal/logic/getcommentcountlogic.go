package logic

import (
	"context"
	"mini-tiktok/service/favorite/pb/favorite"

	"mini-tiktok/service/favorite/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentCountLogic {
	return &GetCommentCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentCountLogic) GetCommentCount(in *favorite.GetCommentCountRequest) (*favorite.GetCommentCountResponse, error) {
	count, err := l.svcCtx.CommentModel.CountByVideoId(uint(in.VideoId))

	return &favorite.GetCommentCountResponse{
		Count: uint64(count),
	}, err
}
