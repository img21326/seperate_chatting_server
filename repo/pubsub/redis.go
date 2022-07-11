package pubsub

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type PubSubRepo struct {
	Redis *redis.Client
}

func NewPubSubRepo(redis *redis.Client) PubSubRepoInterface {
	return &PubSubRepo{
		Redis: redis,
	}
}

func (repo *PubSubRepo) Sub(ctx context.Context, topic string) SubscribeInterface {
	return repo.Redis.Subscribe(ctx, topic)
}

func (repo *PubSubRepo) Pub(ctx context.Context, topic string, message []byte) error {
	return repo.Redis.Publish(ctx, topic, message).Err()
}
