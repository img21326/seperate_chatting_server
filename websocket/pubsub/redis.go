package pubsub

import (
	chatsystem "chat_system"
	"context"

	redis "github.com/redis/go-redis/v9"
)

type RedisPubSub struct {
	redis *redis.Client

	subChan chan []byte
}

func NewRedisPubSub(ctx context.Context, redis *redis.Client) PubSubInterface {

	subChan := make(chan []byte, 1024)
	sub := redis.Subscribe(ctx, chatsystem.SUBSCRIBE_CHANNEL)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(subChan)
				return
			default:
				msg, err := sub.ReceiveMessage(ctx)
				if err != nil {
					panic(err)
				}
				subChan <- []byte(msg.Payload)
			}
		}
	}()

	return &RedisPubSub{
		redis: redis,

		subChan: subChan,
	}
}

func (r *RedisPubSub) Publish(msg []byte) error {
	ctx := context.Background()
	return r.redis.Publish(ctx, chatsystem.SUBSCRIBE_CHANNEL, msg).Err()
}

func (r *RedisPubSub) Subscribe() []byte {
	subMsg := <-r.subChan
	return subMsg
}
