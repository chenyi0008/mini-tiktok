package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"mini-tiktok/service/core/internal/svc"
	"mini-tiktok/service/core/internal/types"
	"mini-tiktok/service/core/models"
	"mini-tiktok/service/favorite/pb/favorite"
)

type GetLikeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLikeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLikeListLogic {
	return &GetLikeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLikeListLogic) GetLikeList(req *types.GetLikeListRequest) (resp *types.GetLikeListResponse, err error) {
	resp = new(types.GetLikeListResponse)
	result, err := l.svcCtx.FavoriteRpc.LikeList(l.ctx, &favorite.LikeListRequest{
		UserId: uint64(req.UserId),
	})
	if err != nil {
		logx.Error(err)
		return
	}

	length := len(result.VideoId)
	idList := make([]int, length)
	idUintList := make([]uint64, length)
	for i, id := range result.VideoId {
		idList[i] = int(id)
		idUintList[i] = id
	}

	videoInfoChan := make(chan *[]*models.VideoModel)
	favoriteCountChan := make(chan *favorite.GetFavoriteCountBatchResponse)
	commentCountChan := make(chan *favorite.GetCommentCountBatchResponse)

	// 视频信息
	go func() {
		batch, err := l.svcCtx.RedisCli.GetVideoInfoBatch(l.ctx, idList)
		// todo mapreduce优化
		for i, video := range batch {
			if video == nil {

				batch[i], err = l.svcCtx.VideoModel.FindOneById(idList[i])
				if err != nil {
					logx.Error(err)
				}
				err := l.svcCtx.RedisCli.SetVideoInfo(l.ctx, batch[i])
				if err != nil {
					logx.Error(err)
				}
			}
		}
		videoInfoChan <- &batch
	}()

	// 点赞数量
	go func() {
		favoriteBatch, err := l.svcCtx.FavoriteRpc.GetFavoriteCountBatch(l.ctx, &favorite.GetFavoriteCountBatchRequest{
			VideoIdList: idUintList,
		})
		if err != nil {
			l.Logger.Error(err)
			favoriteCountChan <- nil
			return
		}
		l.Logger.Info(favoriteBatch)
		favoriteCountChan <- favoriteBatch
	}()

	// 评论数量
	go func() {
		countBatch, err := l.svcCtx.FavoriteRpc.GetCommentCountBatch(l.ctx, &favorite.GetCommentCountBatchRequest{
			VideoIdList: idUintList,
		})
		if err != nil {
			l.Logger.Error(err)
			commentCountChan <- nil
			return
		}
		commentCountChan <- countBatch
	}()

	videoInfoBatch := <-videoInfoChan
	favoriteBatch := <-favoriteCountChan
	countBatch := <-commentCountChan

	for i, item := range *videoInfoBatch {
		videoRes := types.VideoListRes{
			ID: int(item.ID),
			Author: types.AuthorRes{
				ID:   int(item.Author.ID),
				Name: item.Author.Name,
			},
			PlayURL:       item.PlayURL,
			CoverURL:      item.CoverURL,
			CommentCount:  int(countBatch.Count[i]),
			FavoriteCount: int64(favoriteBatch.Count[i]),
		}
		resp.VideoList = append(resp.VideoList, videoRes)
	}
	return
}
