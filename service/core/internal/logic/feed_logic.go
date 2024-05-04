package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/mr"
	"mini-tiktok/service/core/internal/svc"
	"mini-tiktok/service/core/internal/types"
	"mini-tiktok/service/core/models"
	"reflect"
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
	if len(scoreArr) == 0 {
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

	f1 := func(source chan<- interface{}) {
		for _, v := range scoreArr {
			logx.Info(v.Member)
			logx.Info(reflect.TypeOf(v.Member))
			num, err := strconv.Atoi(v.Member.(string))
			if err != nil {
				logx.Error(err)
			}
			source <- num
		}
	}
	f2 := func(item interface{}, writer mr.Writer[*models.VideoModel], cancel func(error)) {
		oneVideo, err := l.svcCtx.VideoModel.FindOneById(item.(int))
		if err != nil {
			cancel(err)
		}
		writer.Write(oneVideo)
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
			PlayURL:  item.PlayURL,
			CoverURL: item.CoverURL,
		}
		resp.VideoList = append(resp.VideoList, videoRes)
	}

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
