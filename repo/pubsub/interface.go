package pubsub

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type SubscribeInterface interface {
	ReceiveMessage(context.Context) (*redis.Message, error)
}

type PubSubRepoInterface interface {
	Sub(ctx context.Context, topic string) SubscribeInterface
	Pub(ctx context.Context, topic string, message []byte) error
}
