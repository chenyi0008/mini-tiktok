package logic

import (
	"context"
	"fmt"
	"mini-tiktok/service/core/internal/svc"
	"mini-tiktok/service/core/internal/types"
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

var t = time.Now()

func caculateTime(s int) {
	t2 := time.Now()
	fmt.Println(s, ":", t2.Sub(t).Milliseconds())
	t = time.Now()
}

func (l *FeedLogic) Feed(req *types.FeedRequest) (resp *types.FeedResponse, err error) {
	resp = new(types.FeedResponse)
	if req.LatestTime == 0 {
		req.LatestTime = uint(time.Now().Unix())
	}
	caculateTime(0)
	caculateTime(1)
	scoreArr, err := l.svcCtx.RedisCli.ZRangeByScore(l.ctx, req.LatestTime)
	if err != nil {
		logx.Error(err)
		return
	}
	caculateTime(2)
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
	caculateTime(3)
	idList := make([]int, length)
	for i, item := range scoreArr {
		num, err := strconv.Atoi(item.Member.(string))
		if err != nil {
			logx.Error(err)
		}
		idList[i] = num
	}
	caculateTime(4)
	batch, err := l.svcCtx.RedisCli.GetVideoInfoBatch(l.ctx, idList)
	caculateTime(44)
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
	caculateTime(5)
	idUintList := make([]uint64, length)
	for i, item := range idList {
		idUintList[i] = uint64(item)
	}
	fmt.Println("idlist:", idUintList)
	favoriteBatch, err := l.svcCtx.FavoriteRpc.GetFavoriteCountBatch(l.ctx, &favorite.GetFavoriteCountBatchRequest{
		VideoIdList: idUintList,
	})
	if err != nil {
		l.Logger.Error(err)
		return
	}
	caculateTime(6)
	countBatch, err := l.svcCtx.FavoriteRpc.GetCommentCountBatch(l.ctx, &favorite.GetCommentCountBatchRequest{
		VideoIdList: idUintList,
	})
	if err != nil {
		l.Logger.Error(err)
		return
	}
	caculateTime(7)

	for i, item := range batch {
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
	caculateTime(8)
	if req.UserId != 0 {
		batch, err := l.svcCtx.FavoriteRpc.IsFavoriteBatch(l.ctx, &favorite.IsFavoriteBatchRequest{
			IsFavoriteList: idUintList,
			UserId:         uint64(req.UserId),
		})
		if err != nil {
			l.Logger.Error(err)
		}
		for i := 0; i < len(resp.VideoList); i++ {
			resp.VideoList[i].IsFavorite = batch.IsFavorite[i]
		}
	}
	caculateTime(9)
	return
}
