package logic

import (
	"context"
	"mini-tiktok/service/core/internal/svc"
	"mini-tiktok/service/core/internal/types"
	"mini-tiktok/service/favorite/pb/favorite"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type VideoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoListLogic {
	return &VideoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoListLogic) VideoList(userId uint) (resp *types.VideoListResponse, err error) {
	resp = new(types.VideoListResponse)
	list, err := l.svcCtx.VideoModel.ListByUserId(userId)
	copier.Copy(&resp.VideoList, &list)
	user, err := l.svcCtx.UserModel.GetById(userId)
	if err != nil {
		return
	}
	for i, video := range resp.VideoList {
		video.Author.ID = int(user.ID)
		video.Author.Name = user.Name
		res, err := l.svcCtx.FavoriteRpc.IsFavorite(l.ctx, &favorite.IsFavoriteRequest{
			UserId:  uint64(user.ID),
			VideoId: uint64(video.ID),
		})
		if err != nil {
			logx.Error(err)
		}
		resp.VideoList[i].IsFavorite = res.IsFavorite

		res2, err := l.svcCtx.FavoriteRpc.GetFavoriteCount(l.ctx, &favorite.GetFavoriteCountRequest{VideoId: uint64(video.ID)})
		if err != nil {
			logx.Error(err)
		}
		resp.VideoList[i].FavoriteCount = int64(res2.Count)

		res3, err := l.svcCtx.FavoriteRpc.GetCommentCount(l.ctx, &favorite.GetCommentCountRequest{VideoId: uint64(video.ID)})
		if err != nil {
			logx.Error(err)
		}
		resp.VideoList[i].CommentCount = int(res3.Count)

	}

	return
}
