package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/go-redis/redis/v8"
	"mini-tiktok/common/consts"
	"mini-tiktok/service/core/define"
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
	fmt.Println(scores.String())
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

func (m *RedisCliModel) SetVideoInfo(ctx context.Context, video *Video) error {
	jsonData, err := json.Marshal(video)
	if err != nil {
		return err
	}
	// 存储对象到 Redis
	key := fmt.Sprintf("%s%d", consts.VideoInfo, video.ID)
	set := m.client.Set(ctx, key, jsonData, 0)
	_, err = set.Result()
	if err != nil {
		return errors.Errorf("err: %s  redis feedback: %s", err.Error(), set.String())
	}
	return nil
}

func (m *RedisCliModel) GetVideoInfo(ctx context.Context, id uint) (*Video, error) {
	key := fmt.Sprintf("%s%d", consts.VideoInfo, id)
	get := m.client.Get(ctx, key)
	result, err := get.Result()
	if err != nil {
		return nil, errors.Errorf("err: %s  redis feedback: %s", err.Error(), get.String())
	}
	video := &Video{}
	err = json.Unmarshal([]byte(result), video)
	if err != nil {
		return nil, err
	}
	return video, nil
}
