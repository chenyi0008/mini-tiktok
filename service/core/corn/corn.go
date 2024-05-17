package corn

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"mini-tiktok/common/consts"
	"mini-tiktok/service/core/internal/svc"
	"strconv"
	"strings"
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
			BatchWriteRedisFavorCountToMySQL(ctx, serverCtx)
			BatchWriteRedisFavorToMySQL(ctx, serverCtx)
		}),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(j.ID())

	// start the scheduler
	s.Start()
}

// BatchWriteRedisFavorCountToMySQL 定时回写点赞数量
func BatchWriteRedisFavorCountToMySQL(ctx context.Context, serverCtx *svc.ServiceContext) error {
	logx.Info("TIMING TASK: BatchWriteRedisFavorCountToMySQL start")
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
			logx.Error("GetVideoFavorCount err: ", err)
			return err
		}
		// 更新数据
		videoIdint, err := strconv.Atoi(videoId)
		if err != nil {
			logx.Error("Atoi err: ", err)
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

// BatchWriteRedisFavorToMySQL 定时回写点赞关系
func BatchWriteRedisFavorToMySQL(ctx context.Context, serverCtx *svc.ServiceContext) error {
	logx.Info("TIMING TASK: BatchWriteRedisFavorToMySQL start")
	// 获取视频点赞tag
	keys, err := serverCtx.RedisCli.GetKeys(ctx, consts.VideoFavorTag)
	if err != nil {
		logx.Error("GetKeys err: ", err)
		return err
	}
	for _, key := range keys {
		videoId, err := parseRedisKeyToVideoId(key)
		if err != nil {
			logx.Error("parseRedisKeyToVideoId err: ", err)
			return err
		}

		tag, err := serverCtx.RedisCli.GetVideoFavorTag(ctx, key)
		if err != nil {
			logx.Error("GetVideoFavorTag err: ", err)
			return err
		}
		for _, t := range tag {
			userId, err := strconv.Atoi(t)
			if err != nil {
				logx.Error("Atoi err: ", err)
				return err
			}
			exist, err := serverCtx.RedisCli.VideoFavorExist(ctx, videoId, userId)
			if err != nil {
				logx.Error("VideoFavorExist err:", err)
			}
			logx.Info("VideoFavorExist videoId: ", videoId, " userId:", userId, "exist:", exist)

			// 如果存在 把数据插入到mysql
			if exist {
				// todo  如果存在 把数据插入到mysql
				logx.Info(" 如果存在 把数据插入到mysql")
				err := serverCtx.FavoriteModel.CreateIgnore(userId, videoId)
				if err != nil {
					logx.Error("CreateIgnore err: ", err)
					return err
				}
			} else { // 不存在 从mysql删除操作
				logx.Info(" 不存在 从mysql删除操作")
				// todo  不存在 从mysql删除操作
				err = serverCtx.FavoriteModel.Del(userId, videoId)
				if err != nil {
					logx.Error("FavoriteModel.Del err:", err)
					return err
				}
			}
		}
	}
	return nil
}

func parseRedisKeyToVideoId(s string) (int, error) {
	// 找到冒号的位置
	index := strings.Index(s, ":")
	if index == -1 {
		return 0, fmt.Errorf("no colon found in string")
	}

	// 提取冒号后的子字符串
	numStr := s[index+1:]

	// 将提取的字符串转换为整数
	num, err := strconv.Atoi(strings.TrimSpace(numStr))
	if err != nil {
		return 0, fmt.Errorf("error converting to integer: %v", err)
	}
	return num, nil
}
