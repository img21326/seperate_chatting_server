package online

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type OnlineRedisRepo struct {
	Redis redis.Cmdable
}

func NewOnlineRedisRepo(redisPipliner redis.Cmdable) OnlineRepoInterface {
	return &OnlineRedisRepo{
		Redis: redisPipliner,
	}
}

func (r *OnlineRedisRepo) Register(ctx context.Context, clientID uint) {
	r.Redis.HSet(ctx, "online", clientID, "1")
}

func (r *OnlineRedisRepo) UnRegister(ctx context.Context, clientID uint) {
	r.Redis.HDel(ctx, "online", strconv.FormatUint(uint64(clientID), 10))
}

func (r *OnlineRedisRepo) CheckUserOnline(ctx context.Context, userId uint) bool {
	redisGet := r.Redis.HGet(ctx, "online", strconv.FormatUint(uint64(userId), 10))
	if redisGet.Err() != nil {
		return false
	} else {
		return true
	}
}
