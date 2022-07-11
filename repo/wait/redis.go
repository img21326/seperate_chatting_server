package wait

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type WaitRedisRepo struct {
	Redis redis.Cmdable
}

func NewRedisWaitRepo(redis redis.Cmdable) WaitRepoInterface {
	return &WaitRedisRepo{
		Redis: redis,
	}
}

func (r *WaitRedisRepo) Add(ctx context.Context, queueName string, clientID uint) {
	r.Redis.RPush(ctx, queueName, clientID)
}

func (r *WaitRedisRepo) Len(ctx context.Context, queueName string) int {
	ret := r.Redis.LLen(ctx, queueName)
	if ret.Err() != nil {
		return 0
	}
	return int(ret.Val())
}

func (r *WaitRedisRepo) Pop(ctx context.Context, queueName string) (clientID uint, err error) {
	ret := r.Redis.LPop(ctx, queueName)
	if ret.Err() != nil {
		log.Printf("[waitRepo] redis get err: %v", ret.Err().Error())
		return 0, ret.Err()
	}
	clientID64, err := ret.Uint64()
	if err != nil {
		log.Printf("[waitRepo] conver to uint err: %v", err)
		return 0, err
	}
	return uint(clientID64), nil
}
