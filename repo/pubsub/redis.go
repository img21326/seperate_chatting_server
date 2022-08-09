package pubsub

import (
	"context"

	"github.com/go-redis/redis/v8"
	errorStruct "github.com/img21326/fb_chat/structure/error"
	"github.com/img21326/fb_chat/structure/pub"
)

type RedisPubSubRepo struct {
	Redis *redis.Client
}

func NewRedisPubSubRepo(redis *redis.Client) PubSubRepoInterface {
	return &RedisPubSubRepo{
		Redis: redis,
	}
}

func (repo *RedisPubSubRepo) Sub(ctx context.Context, topic string) func() ([]byte, error) {
	PubSub := repo.Redis.Subscribe(ctx, topic)
	ReturnChan := make(chan *pub.ReceiveMessage)
	go func(ctx context.Context, PubSub *redis.PubSub, ReturnChan chan *pub.ReceiveMessage) {
		for {
			select {
			case <-ctx.Done():
				RM := &pub.ReceiveMessage{Error: errorStruct.ChannelClosed}
				ReturnChan <- RM
				close(ReturnChan)
				return
			default:
				msg, err := PubSub.ReceiveMessage(ctx)
				RM := &pub.ReceiveMessage{}
				if err != nil {
					RM.Error = err
				} else {
					RM.Payload = []byte(msg.Payload)
					RM.Error = nil
				}
				ReturnChan <- RM
			}
		}
	}(ctx, PubSub, ReturnChan)
	return func() ([]byte, error) {
		rm := <-ReturnChan
		return rm.Payload, rm.Error
	}
}

func (repo *RedisPubSubRepo) Pub(ctx context.Context, topic string, message []byte) error {
	return repo.Redis.Publish(ctx, topic, message).Err()
}
