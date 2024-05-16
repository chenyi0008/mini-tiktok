package corn

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"mini-tiktok/service/core/internal/svc"
	"strconv"
	"time"
)

func CornInit(serverCtx *svc.ServiceContext) {
	s, err := gocron.NewScheduler()
	if err != nil {
		logx.Error(err)
		panic(err)
	}

	ctx := context.Background()
	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(func() {
			BatchWriteRedisDataToMySQL(ctx, serverCtx)
		}),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(j.ID())

	// start the scheduler
	s.Start()
}

func BatchWriteRedisDataToMySQL(ctx context.Context, serverCtx *svc.ServiceContext) error {
	logx.Info("TIMING TASK: BatchWriteRedisDataToMySQL start")
	favorCountSet, err := serverCtx.RedisCli.GetVideoFavorCountTag(ctx)
	if err != nil {
		logx.Error(err)
		return err
	}
	logx.Info("TIMING TASK: FavorCountSet:", favorCountSet)
	for _, videoId := range favorCountSet {
		// 查询数据
		logx.Info("TIMING TASK: videoId:", videoId)
		count, err := serverCtx.RedisCli.GetVideoFavorCount(ctx, videoId)
		if err != nil {
			logx.Error("GetVideoFavorCount err:", err)
			return err
		}
		// 更新数据
		videoIdint, err := strconv.Atoi(videoId)
		if err != nil {
			logx.Error("Atoi err:", err)
			return err
		}
		_, err = serverCtx.VideoModel.UpdateVideoFavorCount(videoIdint, count)
		if err != nil {
			logx.Error("UpdateVideoFavorCount err:", err)
			return err
		}
	}
	logx.Info("TIMING TASK: BatchWriteRedisDataToMySQL Result: Completed ", len(favorCountSet), " task")
	return nil
}
