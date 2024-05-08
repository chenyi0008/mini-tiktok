package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"mini-tiktok/common/consts"
	"mini-tiktok/service/core/define"
	"time"
)

func NewRedisCli(c *redis.Client) *RedisCliModel {
	return &RedisCliModel{
		c,
	}
}

type (
	RedisCliModel struct {
		client *redis.Client
	}
)

func (m *RedisCliModel) ZRangeByScore(ctx context.Context, timestamp uint) ([]redis.Z, error) {
	scores := m.client.ZRangeArgsWithScores(ctx, redis.ZRangeArgs{
		Key:     consts.VideoSortSet,
		Start:   timestamp,
		Stop:    "+inf",
		ByScore: true,
		Count:   int64(define.N),
	})
	result, err := scores.Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *RedisCliModel) ZAddToVideoSortedSet(ctx context.Context, list []VideoModel) error {
	zSetList := make([]*redis.Z, 0)
	for _, item := range list {
		zSetList = append(zSetList, &redis.Z{
			Member: item.ID,
			Score:  float64(item.CreatedAt.Unix()),
		})
	}
	_, err := m.client.ZAdd(ctx, consts.VideoSortSet, zSetList...).Result()
	return err
}

func (m *RedisCliModel) SetVideoInfo(ctx context.Context, video *VideoModel) error {
	jsonData, err := json.Marshal(video)
	if err != nil {
		return err
	}
	// 存储对象到 Redis
	key := fmt.Sprintf("%s%d", consts.VideoInfo, video.ID)

	//生成随机数 防止缓存雪崩
	rand.Seed(time.Now().UnixNano())

	// 生成随机整数
	randomInt := rand.Intn(100)
	set := m.client.Set(ctx, key, jsonData, time.Hour+time.Duration(randomInt)*time.Minute)
	_, err = set.Result()
	if err != nil {
		return errors.Errorf("err: %s  redis feedback: %s", err.Error(), set.String())
	}
	return nil
}

func (m *RedisCliModel) GetVideoInfo(ctx context.Context, id int) (*VideoModel, error) {
	key := fmt.Sprintf("%s%d", consts.VideoInfo, id)
	get := m.client.Get(ctx, key)
	result, err := get.Result()
	if err != nil && err != redis.Nil {
		return nil, errors.Errorf("err: %s  redis feedback: %s", err.Error(), get.String())
	} else if err == redis.Nil {
		return nil, redis.Nil
	}
	video := &VideoModel{}
	err = json.Unmarshal([]byte(result), video)
	if err != nil {
		return nil, err
	}
	return video, err
}

func (m *RedisCliModel) GetVideoInfoBatch(ctx context.Context, ids []int) ([]*VideoModel, error) {
	pipeline := m.client.Pipeline()
	for _, id := range ids {
		key := fmt.Sprintf("%s%d", consts.VideoInfo, id)
		pipeline.Get(ctx, key)
	}
	exec, err := pipeline.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, errors.Errorf("err: %s  redis feedback: %s", err.Error())
	}
	videoList := make([]*VideoModel, 0)
	for _, cmder := range exec {
		cmd := cmder.(*redis.StringCmd)
		result, err := cmd.Result()
		if err != nil && err != redis.Nil {
			return nil, errors.Errorf("err: %s  redis feedback: %s", err.Error())
		}
		var video *VideoModel
		if err == redis.Nil {
			video = nil
		} else {
			video = &VideoModel{}
			err = json.Unmarshal([]byte(result), video)
			if err != nil {
				return nil, err
			}
		}
		videoList = append(videoList, video)

	}
	return videoList, err
}
