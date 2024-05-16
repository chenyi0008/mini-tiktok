package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm/utils"
	"math/rand"
	"mini-tiktok/common/consts"
	"mini-tiktok/service/core/define"
	"strconv"
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
		Start:   "-inf",
		Stop:    timestamp,
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

func (m *RedisCliModel) DelVideoFavorCount(ctx context.Context) error {
	err := m.client.Del(ctx, "video_count:video_favorite_count").Err()
	return err
}

func (m *RedisCliModel) DelVideCommentCount(ctx context.Context) error {
	err := m.client.Del(ctx, "video_count:video_comment_count").Err()
	return err
}

func (m *RedisCliModel) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	var keys []string

	iter := m.client.Scan(ctx, 0, pattern+"*", 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

// GetVideoFavorCountTag 获取视频点赞数量tag集合下的数据
func (m *RedisCliModel) GetVideoFavorCountTag(ctx context.Context) ([]string, error) {
	result, err := m.client.SMembers(ctx, consts.VideoFavoriteCountTag).Result()
	return result, err
}

func (m *RedisCliModel) GetVideoFavorCount(ctx context.Context, id interface{}) (int, error) {
	result, err := m.client.HGet(ctx, consts.VideoFavoriteCount, utils.ToString(id)).Result()
	if err != nil {
		return 0, err
	}
	atoi, err := strconv.Atoi(result)
	if err != nil {
		return 0, err
	}
	return atoi, err
}

//// todo
//func (m *RedisCliModel) ScheduleFavorCountDataWrite(ctx context.Context) error {
//	// 获取tag里的set
//	result, err := m.client.SMembers(ctx, consts.VideoFavoriteCountTag).Result()
//	if err != nil {
//		logx.Error(err)
//		return err
//	}
//	for _, s := range result {
//		// 从redis读数据
//		count, err := m.GetVideoFavorCount(ctx, s)
//		if err != nil {
//			logx.Error(err)
//			return err
//		}
//		// 存入mysql
//
//	}
//
//}
