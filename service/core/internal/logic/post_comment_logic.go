package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"mini-tiktok/service/core/internal/svc"
	"mini-tiktok/service/core/internal/types"
	"mini-tiktok/service/favorite/pb/favorite"
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
	comment, err := l.svcCtx.FavoriteRpc.PostComment(l.ctx, &favorite.PostCommentRequest{
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
