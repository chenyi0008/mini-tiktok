package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"mini-tiktok/core/helper"
	"mini-tiktok/core/internal/svc"
	"mini-tiktok/core/internal/types"
	"mini-tiktok/core/pb/favorite"
)

type PostCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostCommentLogic {
	return &PostCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostCommentLogic) PostComment(req *types.PostCommentRequest) (resp *types.PostCommentResponse, err error) {
	resp = new(types.PostCommentResponse)
	comment, err := helper.FavoriteClient.PostComment(context.Background(), &favorite.PostCommentRequest{
		UserId:  uint64(req.UserId),
		VideoId: uint64(req.VideoId),
		Content: req.CommentText,
	})
	if err != nil {
		return
	}
	resp.Comment.Id = uint(comment.ContentId)
	resp.Comment.CreatedAt = comment.CreatedAt
	resp.Comment.Content = req.CommentText
	resp.Comment.User.Name = req.UserName
	resp.Comment.User.Id = req.UserId
	return
}
