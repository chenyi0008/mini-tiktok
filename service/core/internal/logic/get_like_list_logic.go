package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
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
		//使用mapreduce并发查询
		f1 := func(source chan<- interface{}) {
			for i, v := range batch {
				item := &pair{
					Idx:   i,
					Value: v,
				}
				source <- item
			}
		}
		f2 := func(item interface{}, writer mr.Writer[*pair], cancel func(error)) {
			video := item.(*pair)
			i := video.Idx
			if video.Value == nil {
				value, err := l.svcCtx.VideoModel.FindOneById(idList[i])
				video.Value = value
				if err != nil {
					logx.Error(err)
				}
				err = l.svcCtx.RedisCli.SetVideoInfo(l.ctx, video.Value)
				if err != nil {
					logx.Error(err)
				}
			}
			writer.Write(video)
		}
		f3 := func(pipe <-chan *pair, writer mr.Writer[[]*models.VideoModel], cancel func(error)) {
			lists := make([]*models.VideoModel, length)
			for p := range pipe {
				lists[p.Idx] = p.Value
			}
			writer.Write(lists)
		}
		list, err := mr.MapReduce(f1, f2, f3)
		if err != nil {
			logx.Error(err)
		}

		videoInfoChan <- &list
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
			IsFavorite:    true,
		}
		resp.VideoList = append(resp.VideoList, videoRes)
	}
	return
}
