package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"math/rand"
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
	result, err := m.client.SAdd(ctx, fmt.Sprintf("%s%d", consts.VideoFavor, userId), videoId).Result()
	return int(result), err
}

// GetFavorSet 获取用户关注的
func (m *DefaultRedisCliModel) GetFavorSet(ctx context.Context, userId uint64) ([]uint64, error) {
	res := make([]uint64, 0)
	strArr, err := m.client.SMembers(ctx, fmt.Sprintf("%s%d", consts.VideoFavor, userId)).Result()
	if err != nil {
		return res, err
	}
	for _, s := range strArr {
		atoi, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}
		res = append(res, atoi)
	}
	return res, err
}

// AddFavorSet 添加用户关注的
func (m *DefaultRedisCliModel) AddFavorSet(ctx context.Context, userId uint64, videoList []uint64) error {
	// 构造命令参数
	args := make([]interface{}, len(videoList))
	for i, videoID := range videoList {
		args[i] = videoID
	}

	// 执行批量添加命令
	_, err := m.client.SAdd(ctx, fmt.Sprintf("%s%d", consts.VideoFavor, userId), args...).Result()
	if err != nil {
		return err
	}
	return nil
}

// CancelFavor 取消点赞
func (m *DefaultRedisCliModel) CancelFavor(ctx context.Context, videoId, userId uint64) (int, error) {
	result, err := m.client.SRem(ctx, fmt.Sprintf("%s%d", consts.VideoFavor, userId), videoId).Result()
	return int(result), err
}

// GetIsFavorite 获取是否点赞信息
func (m *DefaultRedisCliModel) GetIsFavorite(ctx context.Context, videoId, userId uint64) (bool, error) {
	setKey := fmt.Sprintf("%s%d", consts.VideoFavor, userId)
	exists, err := m.client.SIsMember(ctx, setKey, videoId).Result()
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

// GetIsFavoriteBatch 批量获取是否点赞信息
func (m *DefaultRedisCliModel) GetIsFavoriteBatch(ctx context.Context, videoId []uint64, userId uint64) ([]bool, error) {
	pipe := m.client.Pipeline()
	length := len(videoId)
	result := make([]bool, length)
	for _, value := range videoId {
		key := fmt.Sprintf("%s%d", consts.VideoFavor, userId)
		pipe.SIsMember(ctx, key, value)
	}
	exec, err := pipe.Exec(ctx)
	if err != nil {
		return result, err
	}

	// 获取每个命令的结果
	for i, cmder := range exec {
		key := fmt.Sprintf("%s%d", consts.VideoFavor, userId)
		cmd, ok := cmder.(*redis.BoolCmd)
		if ok {
			exists, err := cmd.Result()
			str := pipe.SIsMember(ctx, key, videoId[i]).String()
			logx.Info(str)
			if err != nil {
				return result, err
			}
			result[i] = exists
		}
	}
	return result, nil
}

// AddFavorite 添加点赞信息
func (m *DefaultRedisCliModel) AddFavorite(ctx context.Context, videoId, userId uint64) error {
	setKey := fmt.Sprintf("%s%d", consts.VideoFavor, userId)
	_, err := m.client.SAdd(ctx, setKey, videoId).Result()
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

// IncrVideoFavoriteCount 视频点赞数量自增
func (m *DefaultRedisCliModel) IncrVideoFavoriteCount(ctx context.Context, videoId uint64) (int64, error) {
	id := strconv.Itoa(int(videoId))
	result, err := m.client.HIncrBy(ctx, consts.VideoFavoriteCount, id, 1).Result()
	return result, err
}

// DecVideoFavoriteCount 视频点赞数量自减
func (m *DefaultRedisCliModel) DecVideoFavoriteCount(ctx context.Context, videoId uint64) (int64, error) {
	id := strconv.Itoa(int(videoId))
	result, err := m.client.HIncrBy(ctx, consts.VideoFavoriteCount, id, -1).Result()
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

// GetVideoFavoriteCountTag 获取点赞的tag
func (m *DefaultRedisCliModel) GetVideoFavoriteCountTag(ctx context.Context) ([]string, error) {
	set, err := m.client.SMembers(ctx, consts.VideoFavoriteCountTag).Result()
	list := make([]string, 0)
	if err != nil && err != redis.Nil {
		return nil, errors.Errorf("err: %s  redis feedback: %s", err.Error())
	} else if err == redis.Nil {
		return list, nil
	}
	for _, item := range set {
		list = append(list, item)
	}
	return list, err
}

// SetVideoFavoriteCountTag 设置点赞数量的tag
func (m *DefaultRedisCliModel) SetVideoFavoriteCountTag(ctx context.Context, id uint64) error {
	_, err := m.client.SAdd(ctx, consts.VideoFavoriteCountTag, id).Result()
	if err != nil {
		return err
	}
	return nil
}

// SetVideoFavorTag 设置点赞关系的tag
func (m *DefaultRedisCliModel) SetVideoFavorTag(ctx context.Context, videoId, userId uint64) error {
	key := fmt.Sprintf("%s%d", consts.VideoFavorTag, videoId)
	_, err := m.client.SAdd(ctx, key, userId).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetVideoFavorCountTag 获取点赞关系的tag
func (m *DefaultRedisCliModel) GetVideoFavorCountTag(ctx context.Context, videoId int) ([]string, error) {
	key := fmt.Sprintf("%s%d", consts.VideoFavorTag, videoId)
	set, err := m.client.SMembers(ctx, key).Result()
	list := make([]string, 0)
	if err != nil && err != redis.Nil {
		return nil, errors.Errorf("err: %s  redis feedback: %s", err.Error())
	} else if err == redis.Nil {
		return list, nil
	}
	for _, item := range set {
		list = append(list, item)
	}
	return list, err
}

// GetComment 获取评论
func (m *DefaultRedisCliModel) GetComment(ctx context.Context, videoId uint64) ([]Comment, error) {
	key := fmt.Sprintf("%s%d", consts.VideoComment, videoId)
	set, err := m.client.ZRandMember(ctx, key, 100, false).Result()
	list := make([]Comment, 0)
	if err != nil && err != redis.Nil {
		return nil, errors.Errorf("err: %s  redis feedback: %s", err.Error())
	} else if len(set) == 1 && set[0] == "tag" {
		//防止缓存击穿 处理tag
		return list, err
	} else if len(set) == 0 {
		// 未命中
		_, err := m.client.ZAdd(ctx, key, &redis.Z{
			Score:  0,
			Member: "tag",
		}).Result()
		if err != nil {
			logx.Error("client.SAdd err: ", err)
			return nil, err
		}
		return list, redis.Nil
	}

	// 处理tag标签
	if set[len(set)-1] == "tag" {
		set = set[:len(set)-1]
	} else {
		_, err := m.client.ZAdd(ctx, key, &redis.Z{
			Score:  0,
			Member: "tag",
		}).Result()
		if err != nil {
			logx.Error("client.SAdd err: ", err)
			return nil, err
		}
	}
	for _, item := range set {
		comment := Comment{}
		err := json.Unmarshal([]byte(item), &comment)
		if err != nil {
			logx.Error("json.Unmarshal err:", err)
			return nil, err
		}
		list = append(list, comment)
	}
	return list, err
}

// AddComment 添加评论
func (m *DefaultRedisCliModel) AddComment(ctx context.Context, videoId uint64, commentList []Comment) error {

	key := fmt.Sprintf("%s%d", consts.VideoComment, videoId)
	mList := make([]*redis.Z, 0)
	for _, comment := range commentList {
		marshal, err := json.Marshal(comment)
		if err != nil {
			logx.Error("json.Marshal err:", err)
			return err
		}
		mList = append(mList, &redis.Z{
			Score:  float64(comment.CreatedAt.Unix()),
			Member: string(marshal),
		})
	}

	_, err := m.client.ZAdd(ctx, key, mList...).Result()
	if err != nil {
		logx.Error("client.ZAdd err: ", err)
		return err
	}

	//添加过期时间
	_, err = m.client.Expire(ctx, key, hour()).Result()
	if err != nil {
		log.Fatalf("Could not set expiration for key: %v", err)
	}

	return nil
}

func (m *DefaultRedisCliModel) DelComment(ctx context.Context, videoId uint64) error {
	key := fmt.Sprintf("%s%d", consts.VideoComment, videoId)
	_, err := m.client.Del(ctx, key).Result()
	return err
}

func hour() time.Duration {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(2) + 1
	return time.Duration(randomNumber) * time.Hour
}
