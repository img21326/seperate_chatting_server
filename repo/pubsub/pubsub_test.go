package pubsub

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/img21326/fb_chat/helper"
	errorStruct "github.com/img21326/fb_chat/structure/error"
	"github.com/stretchr/testify/assert"
)

func getRedis() *redis.Client {
	var redisClient = &redis.Options{
		Addr:     "139.162.125.28:6379",
		Password: "",
		DB:       5,
	}
	return redis.NewClient(redisClient)
}

func TestPubSub(t *testing.T) {
	redis := getRedis()
	repo := PubSubRepo{Redis: redis}

	ctx := context.Background()
	topic := helper.RandString(5)
	subFunc := repo.Sub(ctx, topic)
	repo.Pub(ctx, topic, []byte("test"))

	msg, err := subFunc()

	assert.Nil(t, err)
	assert.Equal(t, string(msg[:]), "test")
}

func TestSubShutdown(t *testing.T) {
	redis := getRedis()
	repo := PubSubRepo{Redis: redis}

	ctx, cancel := context.WithCancel(context.Background())
	topic := helper.RandString(5)
	subFunc := repo.Sub(ctx, topic)
	cancel()
	msg, err := subFunc()

	assert.Nil(t, msg)
	assert.Equal(t, err, errorStruct.ChannelClosed)
}
