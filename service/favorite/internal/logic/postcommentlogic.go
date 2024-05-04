package logic

import (
	"context"
	"mini-tiktok/service/favorite/pb/favorite"

	"mini-tiktok/service/favorite/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostCommentLogic {
	return &PostCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PostCommentLogic) PostComment(in *favorite.PostCommentRequest) (*favorite.PostCommentResponse, error) {
	resp := new(favorite.PostCommentResponse)
	comment, err := l.svcCtx.CommentModel.Create(uint(in.UserId), uint(in.VideoId), in.Content)
	if err != nil {
		resp.Message = "评论失败"
		resp.Code = 1
		return resp, nil
	}

	resp.Message = "评论成功"
	resp.ContentId = uint64(comment.ID)
	resp.CreatedAt = comment.CreatedAt.Format("2006-01-02 15:04:05")
	return resp, nil
}
