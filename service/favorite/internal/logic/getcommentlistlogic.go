package logic

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"golang.org/x/exp/slices"
	"mini-tiktok/service/favorite/pb/favorite"

	"mini-tiktok/service/favorite/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentListLogic {
	return &GetCommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentListLogic) GetCommentList(in *favorite.GetCommentRequest) (*favorite.GetCommentResponse, error) {
	resp := new(favorite.GetCommentResponse)
	commentList, err := l.svcCtx.RedisCli.GetComment(l.ctx, in.VideoId)
	if err != nil && err != redis.Nil {
		logx.Error("RedisCli.GetComment err: ", err)
		return nil, err
	} else if err == redis.Nil {
		//未命中 第一次加载
		logx.Info("GetCommentList 未命中：", in.VideoId)
		commentList, err = l.svcCtx.CommentModel.GetByVideoId(uint(in.VideoId))
		if err != nil {
			logx.Error("CommentModel.GetByVideoId err:", err)
			return nil, err
		}
		if len(commentList) == 0 {

		} else {
			err = l.svcCtx.RedisCli.AddComment(l.ctx, in.VideoId, commentList)
			if err != nil {
				logx.Error("RedisCli.AddComment err: ", err)
				return nil, err
			}
		}
		slices.Reverse(commentList)
	}
	copier.Copy(&resp.CommentList, &commentList)

	return resp, nil
}
