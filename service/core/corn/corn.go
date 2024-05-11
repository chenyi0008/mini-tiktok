package corn

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"mini-tiktok/service/core/internal/svc"
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

func BatchWriteRedisDataToMySQL(ctx context.Context, serverCtx *svc.ServiceContext) {
	logx.Info("timing task: BatchWriteRedisDataToMySQL")
	info, err := serverCtx.RedisCli.GetVideoInfo(ctx, 13)
	if err != nil {
		return
	}
	fmt.Println(info)
}
