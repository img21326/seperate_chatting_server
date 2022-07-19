package wait

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return s
}

func getRedis() *redis.Client {
	redisServer := mockRedis()
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	return redisClient
}

func TestAdd(t *testing.T) {
	redis := getRedis()
	waitRepo := WaitRedisRepo{
		Redis: redis,
	}
	ctx := context.Background()
	r := redis.LLen(ctx, "test")
	assert.Equal(t, int(r.Val()), 0)
	waitRepo.Add(ctx, "test", 1)
	r = redis.LLen(ctx, "test")
	assert.Equal(t, int(r.Val()), 1)
}

func TestLen(t *testing.T) {
	redis := getRedis()
	waitRepo := WaitRedisRepo{
		Redis: redis,
	}
	ctx := context.Background()
	r := waitRepo.Len(ctx, "test")
	assert.Equal(t, r, 0)
	waitRepo.Add(ctx, "test", 1)
	r = waitRepo.Len(ctx, "test")
	assert.Equal(t, r, 1)
}

func TestPop(t *testing.T) {
	redis := getRedis()
	waitRepo := WaitRedisRepo{
		Redis: redis,
	}
	ctx := context.Background()
	waitRepo.Add(ctx, "test", 1)
	r, err := waitRepo.Pop(ctx, "test")
	assert.Equal(t, r, uint(1))
	assert.Nil(t, err)
}

func TestPopErr(t *testing.T) {
	redis := getRedis()
	waitRepo := WaitRedisRepo{
		Redis: redis,
	}
	ctx := context.Background()
	r, err := waitRepo.Pop(ctx, "test")
	t.Logf("pop r: %+v, err: %+v", r, err)
	assert.NotNil(t, err)
}
