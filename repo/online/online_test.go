package online

import (
	"context"
	"sync"
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

func TestLocalRegister(t *testing.T) {
	onlineRepo := &LocalOnlineRepo{
		ClientMap: make(map[uint]bool),
		lock:      &sync.RWMutex{},
	}

	ctx := context.Background()
	onlineRepo.Register(ctx, 1)

	assert.Equal(t, onlineRepo.ClientMap[1], true)
}

func TestLocalUnRegister(t *testing.T) {
	onlineRepo := &LocalOnlineRepo{
		ClientMap: make(map[uint]bool),
		lock:      &sync.RWMutex{},
	}

	ctx := context.Background()
	onlineRepo.Register(ctx, 1)
	onlineRepo.UnRegister(ctx, 1)
	assert.Equal(t, onlineRepo.ClientMap[1], false)
}

func TestLocalCheckUserOnline(t *testing.T) {
	onlineRepo := &LocalOnlineRepo{
		ClientMap: make(map[uint]bool),
		lock:      &sync.RWMutex{},
	}

	ctx := context.Background()
	assert.Equal(t, onlineRepo.CheckUserOnline(ctx, 1), false)
	onlineRepo.Register(ctx, 1)
	assert.Equal(t, onlineRepo.CheckUserOnline(ctx, 1), true)
	onlineRepo.UnRegister(ctx, 1)
	assert.Equal(t, onlineRepo.CheckUserOnline(ctx, 1), false)
}

func TestRedisRegister(t *testing.T) {
	redis := getRedis()
	onlineRepo := OnlineRedisRepo{Redis: redis}

	ctx := context.Background()
	onlineRepo.Register(ctx, 1)

	l := redis.HGet(ctx, "online", "1")
	assert.Equal(t, l.Val(), "1")
}

func TestRedisUnRegister(t *testing.T) {
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

func TestRedisCheckUserOnline(t *testing.T) {
	redis := getRedis()
	onlineRepo := OnlineRedisRepo{Redis: redis}

	ctx := context.Background()
	onlineRepo.Register(ctx, 1)
	assert.Equal(t, onlineRepo.CheckUserOnline(ctx, 1), true)
	onlineRepo.UnRegister(ctx, 1)
	assert.Equal(t, onlineRepo.CheckUserOnline(ctx, 1), false)
}
