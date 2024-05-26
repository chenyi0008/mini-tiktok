package logic

import (
	"context"
	"mini-tiktok/common/consts"
	"mini-tiktok/service/favorite/internal/svc"
	"mini-tiktok/service/favorite/pb/favorite"

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
		logx.Error("RedisCli.AddComment err: ", err)
		resp.Message = "评论失败"
		resp.Code = consts.FAILED
		return resp, nil
	}
	//延迟双删
	err = l.svcCtx.RedisCli.DelComment(l.ctx, in.VideoId)
	if err != nil {
		logx.Error(err)
	}

	resp.Message = "评论成功"
	resp.ContentId = uint64(comment.ID)
	resp.CreatedAt = comment.CreatedAt.Format("2006-01-02 15:04:05")
	return resp, nil
}
