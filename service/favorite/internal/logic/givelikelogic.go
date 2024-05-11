package logic

import (
	"context"
	"mini-tiktok/common/consts"
	"mini-tiktok/service/favorite/pb/favorite"
	"sync"

	"mini-tiktok/service/favorite/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GiveLikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGiveLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GiveLikeLogic {
	return &GiveLikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GiveLikeLogic) GiveLike(in *favorite.GiveLikeRequest) (*favorite.Response, error) {
	// 判断是否已点赞
	isFavorite, err := l.svcCtx.RedisCli.GetIsFavorite(l.ctx, in.VideoId, in.UserId)
	if err != nil {
		logx.Error(err)
		return &favorite.Response{
			Code:    consts.FAILED,
			Message: "点赞失败",
		}, err
	}
	if isFavorite {
		return &favorite.Response{
			Code:    consts.FAILED,
			Message: "点赞失败",
		}, err
	}

	// 使用 WaitGroup 等待所有的 goroutine 完成
	var wg sync.WaitGroup
	// 定义一个错误通道，用于收集所有 goroutine 的错误
	errChan := make(chan error, 4)

	// 开启 goroutine 执行每个 Redis 操作
	wg.Add(4)

	// 添加点赞关系
	go func() {
		defer wg.Done()
		_, err := l.svcCtx.RedisCli.SetFavor(l.ctx, in.VideoId, in.UserId)
		errChan <- err
	}()

	// 添加点赞数量
	go func() {
		defer wg.Done()
		_, err := l.svcCtx.RedisCli.IncrVideoFavoriteCount(l.ctx, in.VideoId)
		errChan <- err
	}()

	// 添加点赞关系tag
	go func() {
		defer wg.Done()
		err := l.svcCtx.RedisCli.SetVideoFavorTag(l.ctx, in.VideoId, in.UserId)
		errChan <- err
	}()

	// 添加点赞数量tag
	go func() {
		defer wg.Done()
		err := l.svcCtx.RedisCli.SetVideoFavoriteCountTag(l.ctx, in.VideoId)
		errChan <- err
	}()

	// 等待所有 goroutine 完成
	wg.Wait()
	// 关闭错误通道，确保在读取完所有错误后，不会阻塞主线程
	close(errChan)

	// 遍历错误通道，检查是否有错误发生
	for err := range errChan {
		if err != nil {
			logx.Error(err)
			return &favorite.Response{
				Code:    consts.FAILED,
				Message: "点赞失败",
			}, err
		}
	}

	// 如果没有错误发生，返回成功响应
	return &favorite.Response{
		Code:    consts.SUCCEED,
		Message: "点赞成功",
	}, nil

}
