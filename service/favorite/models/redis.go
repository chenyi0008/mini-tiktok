package models

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"mini-tiktok/common/consts"
	"reflect"
	"strconv"
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

// GetIsFavoriteBatch 批量获取点赞信息
func (m *DefaultRedisCliModel) GetIsFavoriteBatch(ctx context.Context, videoId []uint64, userId uint64) ([]bool, error) {
	pipe := m.client.Pipeline()
	length := len(videoId)
	result := make([]bool, length)
	for _, value := range videoId {
		key := fmt.Sprintf("%s%d", consts.VideoFavor, value)
		pipe.SIsMember(ctx, key, userId)
	}
	exec, err := pipe.Exec(ctx)
	if err != nil {
		return result, err
	}

	// 获取每个命令的结果
	for i, cmder := range exec {
		key := fmt.Sprintf("%s%d", consts.VideoFavor, videoId[i])
		cmd, ok := cmder.(*redis.BoolCmd)
		if ok {
			exists, err := cmd.Result()
			str := pipe.SIsMember(ctx, key, userId).String()
			logx.Info(str)
			if err != nil {
				return result, err
			}
			result[i] = exists
		}
	}
	return result, nil
}

func (m *DefaultRedisCliModel) AddFavorite(ctx context.Context, videoId, userId uint64) error {
	setKey := fmt.Sprintf("%s%d", consts.VideoFavor, videoId)
	_, err := m.client.SAdd(ctx, setKey, userId).Result()
	return err
}

// GetFavoriteCountBatch 批量获取点赞数量
func (m *DefaultRedisCliModel) GetFavoriteCountBatch(ctx context.Context, videoId []uint64) ([]int, []int, error) {
	missedRecordIdx := make([]int, 0)
	pipe := m.client.Pipeline()
	length := len(videoId)
	result := make([]int, length)
	for _, value := range videoId {
		pipe.HGet(ctx, consts.VideoFavoriteCount, strconv.FormatUint(value, 10))
	}
	exec, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, nil, err
	}
	// 获取每个命令的结果
	for i, cmder := range exec {
		cmd, ok := cmder.(*redis.StringCmd)
		if ok {
			res, err := cmd.Result()
			if err != nil && err != redis.Nil {
				panic(err)
				fmt.Println(reflect.TypeOf(err))
				return nil, nil, err
			}
			if err == redis.Nil {
				missedRecordIdx = append(missedRecordIdx, i)
				result[i] = -1
			} else {
				atoi, err := strconv.Atoi(res)
				if err != nil {
					return nil, nil, err
				}
				result[i] = atoi
			}
		}
	}
	return result, missedRecordIdx, nil
}

func (m *DefaultRedisCliModel) SetVideoFavoriteCount(ctx context.Context, videoId, count int) (int64, error) {
	result, err := m.client.HSet(ctx, consts.VideoFavoriteCount, videoId, count).Result()
	return result, err
}

// GetCommentCountBatch 批量获取评论数量
func (m *DefaultRedisCliModel) GetCommentCountBatch(ctx context.Context, videoId []uint64) ([]int, []int, error) {
	missedRecordIdx := make([]int, 0)
	pipe := m.client.Pipeline()
	length := len(videoId)
	result := make([]int, length)
	for _, value := range videoId {
		pipe.HGet(ctx, consts.VideoCommentCount, strconv.FormatUint(value, 10))
	}
	exec, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, nil, err
	}
	// 获取每个命令的结果
	for i, cmder := range exec {
		cmd, ok := cmder.(*redis.StringCmd)
		if ok {
			res, err := cmd.Result()
			if err != nil && err != redis.Nil {
				panic(err)
				fmt.Println(reflect.TypeOf(err))
				return nil, nil, err
			}
			if err == redis.Nil {
				missedRecordIdx = append(missedRecordIdx, i)
				result[i] = -1
			} else {
				atoi, err := strconv.Atoi(res)
				if err != nil {
					return nil, nil, err
				}
				result[i] = atoi
			}
		}
	}
	return result, missedRecordIdx, nil
}

func (m *DefaultRedisCliModel) SetVideoCommentCount(ctx context.Context, videoId, count int) (int64, error) {
	result, err := m.client.HSet(ctx, consts.VideoCommentCount, videoId, count).Result()
	return result, err
}
