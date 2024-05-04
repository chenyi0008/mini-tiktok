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

func NewRedisCli(c *redis.Client) *DefaultRedisCliModel {
	return &DefaultRedisCliModel{
		c,
	}
}

type (
	DefaultRedisCliModel struct {
		client *redis.Client
	}
)

func (m *DefaultRedisCliModel) NumOfFavor(ctx context.Context, videoId uint64) (int, error) {
	result, err := m.client.SCard(ctx, fmt.Sprintf("%s%d", consts.VideoFavor, videoId)).Result()
	if err != nil {
		return -1, err
	}
	return int(result), nil
}
func (m *DefaultRedisCliModel) SetFavor(ctx context.Context, videoId, userId uint64) (int, error) {
	result, err := m.client.SAdd(ctx, fmt.Sprintf("%s%d", consts.VideoFavor, videoId), userId).Result()
	return int(result), err
}

func (m *DefaultRedisCliModel) CancelFavor(ctx context.Context, videoId, userId uint64) (int, error) {
	result, err := m.client.SRem(ctx, fmt.Sprintf("%s%d", consts.VideoFavor, videoId), userId).Result()
	return int(result), err
}

func (m *DefaultRedisCliModel) GetIsFavorite(ctx context.Context, videoId, userId uint64) (bool, error) {
	setKey := fmt.Sprintf("%s%d", consts.VideoFavor, videoId)
	exists, err := m.client.SIsMember(ctx, setKey, userId).Result()
	if err != nil {
		return false, err
	}
	// 根据查询结果进行处理
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func (m *DefaultRedisCliModel) AddFavorite(ctx context.Context, videoId, userId uint64) error {
	setKey := fmt.Sprintf("%s%d", consts.VideoFavor, videoId)
	_, err := m.client.SAdd(ctx, setKey, userId).Result()
	return err
}
