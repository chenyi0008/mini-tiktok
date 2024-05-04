package models

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"mini-tiktok/common/consts"
	"sync"
	"time"
)

type RedisHelper struct {
	*redis.Client
}

var redisHelper *RedisHelper

var redisOnce sync.Once

func NewRedisHelper(addr, password string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	},
	)

	redisOnce.Do(func() {
		rdh := new(RedisHelper)
		rdh.Client = rdb
		redisHelper = rdh
	})

	return rdb
}

func InitRedis(addr, password string, db int) *redis.Client {
	rdb := NewRedisHelper(addr, password, db)
	ctx := context.Background()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(err)
		return nil
	}
	return rdb
}

func NewRedisCli(c *redis.Client) *DefaultModel {
	return &DefaultModel{
		c,
	}
}

type (
	DefaultModel struct {
		client *redis.Client
	}
)

func (m *DefaultModel) NumOfFavor(ctx context.Context, VideoId uint64) (int, error) {
	result, err := m.client.SCard(ctx, fmt.Sprintf("%s%d", consts.VideoFavor, VideoId)).Result()
	if err != nil {
		return -1, err
	}
	return int(result), nil
}
func (m *DefaultModel) SetFavor(ctx context.Context, VideoId, UserId uint64) (int, error) {
	result, err := m.client.SAdd(ctx, fmt.Sprintf("%s%d", consts.VideoFavor, VideoId), UserId).Result()
	return int(result), err
}

func (m *DefaultModel) CancelFavor(ctx context.Context, VideoId, UserId uint64) (int, error) {
	result, err := m.client.SRem(ctx, fmt.Sprintf("%s%d", consts.VideoFavor, VideoId), UserId).Result()
	return int(result), err
}
