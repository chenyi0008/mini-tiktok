package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/mr"
	"mini-tiktok/service/core/internal/svc"
	"mini-tiktok/service/core/internal/types"
	"mini-tiktok/service/core/models"
	"mini-tiktok/service/favorite/pb/favorite"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type FeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type pair struct {
	Idx   int
	Value *models.VideoModel
}

func (l *FeedLogic) Feed(req *types.FeedRequest) (resp *types.FeedResponse, err error) {
	resp = new(types.FeedResponse)
	if req.LatestTime == 0 {
		req.LatestTime = uint(time.Now().Unix())
	}
	// 从redis取缓存
	scoreArr, err := l.svcCtx.RedisCli.ZRangeByScore(l.ctx, req.LatestTime)
	if err != nil {
		logx.Error(err)
		return
	}

	// redis没有缓存 从mysql读取
	length := len(scoreArr)
	if length == 0 {
		videoList, err2 := l.svcCtx.VideoModel.ListAll()
		if err2 != nil {
			logx.Error(err2)
			return
		}

		err2 = l.svcCtx.RedisCli.ZAddToVideoSortedSet(l.ctx, videoList)
		if err2 != nil {
			logx.Error(err2)
			return
		}

		scoreArr, err2 = l.svcCtx.RedisCli.ZRangeByScore(l.ctx, req.LatestTime)
		if err2 != nil {
			logx.Error(err2)
			return
		}
	}

	idList := make([]int, length)
	idUintList := make([]uint64, length)
	for i, item := range scoreArr {
		num, err := strconv.Atoi(item.Member.(string))
		if err != nil {
			logx.Error(err)
		}
		idList[i] = num
		idUintList[i] = uint64(num)
	}
	l.Logger.Info(idList)

	videoInfoChan := make(chan *[]*models.VideoModel)
	favoriteCountChan := make(chan *favorite.GetFavoriteCountBatchResponse)
	commentCountChan := make(chan *favorite.GetCommentCountBatchResponse)
	userFavoriteChan := make(chan *favorite.IsFavoriteBatchResponse)

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

	// 是否点赞
	if req.UserId != 0 {
		go func() {
			batch, err := l.svcCtx.FavoriteRpc.IsFavoriteBatch(l.ctx, &favorite.IsFavoriteBatchRequest{
				IsFavoriteList: idUintList,
				UserId:         uint64(req.UserId),
			})
			if err != nil {
				l.Logger.Error(err)
			}
			userFavoriteChan <- batch
		}()
	}

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

	if req.UserId != 0 {
		batch := <-userFavoriteChan
		for i := 0; i < len(resp.VideoList); i++ {
			resp.VideoList[i].IsFavorite = batch.IsFavorite[i]
		}
	}
	return
}
