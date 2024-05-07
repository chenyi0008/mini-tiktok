package logic

import (
	"context"
	"github.com/go-redis/redis/v8"
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

func (l *FeedLogic) Feed(req *types.FeedRequest) (resp *types.FeedResponse, err error) {
	resp = new(types.FeedResponse)
	if req.LatestTime == 0 {
		req.LatestTime = uint(time.Now().Unix())
	}

	scoreArr, err := l.svcCtx.RedisCli.ZRangeByScore(l.ctx, req.LatestTime)
	if err != nil {
		logx.Error(err)
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

	idList := make([]uint64, length)
	for i, item := range scoreArr {
		num, err := strconv.ParseUint(item.Member.(string), 10, 32)
		if err != nil {
			logx.Error(err)
		}
		idList[i] = num
	}

	//使用mapreduce并发查询
	f1 := func(source chan<- interface{}) {
		for _, v := range idList {
			source <- v
		}
	}
	f2 := func(item interface{}, writer mr.Writer[*models.VideoModel], cancel func(error)) {
		id := item.(uint64)
		video, err := l.svcCtx.RedisCli.GetVideoInfo(l.ctx, int(id))
		if err != nil && err != redis.Nil {
			logx.Error(err)
		}
		if video == nil {
			video, err = l.svcCtx.VideoModel.FindOneById(int(id))
			if err != nil {
				logx.Error(err)
				cancel(err)
			}
			err := l.svcCtx.RedisCli.SetVideoInfo(l.ctx, video)
			if err != nil {
				logx.Error(err)
				cancel(err)
			}
		}
		writer.Write(video)
	}
	f3 := func(pipe <-chan *models.VideoModel, writer mr.Writer[[]*models.VideoModel], cancel func(error)) {
		lists := make([]*models.VideoModel, 0)
		for p := range pipe {
			lists = append(lists, p)
		}
		writer.Write(lists)
	}
	list, err := mr.MapReduce(f1, f2, f3)
	for _, item := range list {
		videoRes := types.VideoListRes{
			ID: int(item.ID),
			Author: types.AuthorRes{
				ID:   int(item.Author.ID),
				Name: item.Author.Name,
			},
			PlayURL:       item.PlayURL,
			CoverURL:      item.CoverURL,
			CommentCount:  0, // todo comment count
			FavoriteCount: 0, // todo favorite count
		}
		resp.VideoList = append(resp.VideoList, videoRes)
	}

	if req.UserId != 0 {
		batch, err := l.svcCtx.FavoriteRpc.IsFavoriteBatch(l.ctx, &favorite.IsFavoriteBatchRequest{
			IsFavoriteList: idList,
			UserId:         uint64(req.UserId),
		})
		if err != nil {
			l.Logger.Error(err)
		}

		for i := 0; i < len(resp.VideoList); i++ {
			resp.VideoList[i].IsFavorite = batch.IsFavorite[i]
		}
	}

	favoriteBatch, err := l.svcCtx.FavoriteRpc.GetFavoriteCountBatch(l.ctx, &favorite.GetFavoriteCountBatchRequest{
		VideoIdList: idList,
	})
	countBatch, err := l.svcCtx.FavoriteRpc.GetCommentCountBatch(l.ctx, &favorite.GetCommentCountBatchRequest{
		VideoIdList: idList,
	})
	if err != nil {
		l.Logger.Error(err)
	}
	for i := 0; i < len(resp.VideoList); i++ {
		resp.VideoList[i].FavoriteCount = int64(favoriteBatch.Count[i])
		resp.VideoList[i].CommentCount = int(countBatch.Count[i])
	}
	//	if err != nil {
	//		logx.Error(err)
	//	}

	//videos, _ := l.svcCtx.VideoModel.ListByCreatedAt(int64(req.LatestTime), uint(define.N))
	//copier.Copy(&resp.VideoList, &videos)
	//
	//length := len(videos)
	//if length == 0 {
	//	videos, _ = l.svcCtx.VideoModel.ListByCreatedAt(time.Now().Unix(), uint(define.N))
	//	length = len(videos)
	//}
	//resp.NextTime = uint64(videos[length-1].CreatedAt.Unix())
	//
	//for i, item := range resp.VideoList {
	//	if req.UserId == 0 {
	//		resp.VideoList[i].IsFavorite = false
	//	} else {
	//		res, err := l.svcCtx.FavoriteRpc.IsFavorite(l.ctx, &favorite.IsFavoriteRequest{
	//			UserId:  uint64(req.UserId),
	//			VideoId: uint64(item.ID),
	//		})
	//		if err != nil {
	//			logx.Error(err)
	//		}
	//		resp.VideoList[i].IsFavorite = res.IsFavorite
	//	}
	//
	//	res2, err := l.svcCtx.FavoriteRpc.GetFavoriteCount(l.ctx, &favorite.GetFavoriteCountRequest{VideoId: uint64(item.ID)})
	//	if err != nil {
	//		logx.Error(err)
	//	}
	//	resp.VideoList[i].FavoriteCount = int64(res2.Count)
	//
	//	res3, err := l.svcCtx.FavoriteRpc.GetCommentCount(l.ctx, &favorite.GetCommentCountRequest{VideoId: uint64(item.ID)})
	//	if err != nil {
	//		logx.Error(err)
	//	}
	//	resp.VideoList[i].CommentCount = int(res3.Count)
	//}

	return
}
