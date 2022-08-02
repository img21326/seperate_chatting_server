package pubsub

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/img21326/fb_chat/structure/pub"
)

type PubSubRepo struct {
	Redis *redis.Client
}

func NewPubSubRepo(redis *redis.Client) PubSubRepoInterface {
	return &PubSubRepo{
		Redis: redis,
	}
}

func (repo *PubSubRepo) Sub(ctx context.Context, topic string) func() ([]byte, error) {
	PubSub := repo.Redis.Subscribe(ctx, topic)
	ReturnChan := make(chan *pub.ReceiveMessage)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ReturnChan)
				return
			default:
				msg, err := PubSub.ReceiveMessage(ctx)
				RM := &pub.ReceiveMessage{Payload: []byte(msg.Payload), Error: err}
				ReturnChan <- RM
			}
		}
	}()
	return func() ([]byte, error) {
		rm := <-ReturnChan
		return rm.Payload, rm.Error
	}
}

func (repo *PubSubRepo) Pub(ctx context.Context, topic string, message []byte) error {
	return repo.Redis.Publish(ctx, topic, message).Err()
}
