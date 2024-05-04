package logic

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
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
	list, err := l.svcCtx.CommentModel.GetByVideoId(uint(in.VideoId))
	fmt.Printf("list:%+v\n", list)
	if err != nil {
		resp.Code = 1
		resp.Message = err.Error()
	}
	copier.Copy(&resp.CommentList, &list)

	return resp, nil
}
