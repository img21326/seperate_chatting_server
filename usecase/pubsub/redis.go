package pubsub

import (
	"context"
	"encoding/json"
	"log"

	repo "github.com/img21326/fb_chat/repo/pubsub"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
)

type RedisSubUsecase struct {
	PubSubRepo repo.PubSubRepoInterface
}

func NewRedisSubUsecase(pubSubRepo repo.PubSubRepoInterface) SubMessageUsecaseInterface {
	return &RedisSubUsecase{
		PubSubRepo: pubSubRepo,
	}
}

func (u *RedisSubUsecase) Subscribe(ctx context.Context, topic string) repo.SubscribeInterface {
	subscriber := u.PubSubRepo.Sub(ctx, topic)
	return subscriber
}

func (u *RedisSubUsecase) Publish(ctx context.Context, topic string, message *pubmessage.PublishMessage) {
	log.Printf("[PublishUsecase] send message: %v", message)
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("pub message convert json error: %v", err)
	}
	u.PubSubRepo.Pub(ctx, topic, []byte(jsonMessage))
}
