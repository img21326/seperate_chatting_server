package online

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var redisServer *miniredis.Miniredis

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

func teardown() {
	redisServer.Close()
}

func TestRegister(t *testing.T) {
	redis := getRedis()
	onlineRepo := OnlineRedisRepo{Redis: redis}

	ctx := context.Background()
	onlineRepo.Register(ctx, 1)

	l := redis.HGet(ctx, "online", "1")
	assert.Equal(t, l.Val(), "1")
}

func TestUnRegister(t *testing.T) {
	redis := getRedis()
	onlineRepo := OnlineRedisRepo{Redis: redis}

	ctx := context.Background()
	onlineRepo.Register(ctx, 1)
	l := redis.HLen(ctx, "online")
	assert.Equal(t, int(l.Val()), 1)
	onlineRepo.UnRegister(ctx, 1)
	l = redis.HLen(ctx, "online")
	assert.Equal(t, int(l.Val()), 0)
}

func TestCheckUserOnline(t *testing.T) {
	redis := getRedis()
	onlineRepo := OnlineRedisRepo{Redis: redis}

	ctx := context.Background()
	onlineRepo.Register(ctx, 1)
	assert.Equal(t, onlineRepo.CheckUserOnline(ctx, 1), true)
	onlineRepo.UnRegister(ctx, 1)
	assert.Equal(t, onlineRepo.CheckUserOnline(ctx, 1), false)
}
