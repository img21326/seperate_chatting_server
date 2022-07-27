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

func (u *RedisSubUsecase) Subscribe(ctx context.Context, topic string, processMessage func(*pubmessage.PublishMessage)) {
	log.Printf("[Sub] start subscribe %v", topic)
	subscriber := u.PubSubRepo.Sub(ctx, topic)
	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		log.Printf("[Sub] get message: %v", msg)
		if err != nil {
			log.Printf("[Sub] sub message receive error: %v", err)
		}
		var redisMessage pubmessage.PublishMessage

		if err := json.Unmarshal([]byte(msg.Payload), &redisMessage); err != nil {
			log.Printf("pubsub message json load error: %v", err)
		}
		processMessage(&redisMessage)
	}
}

func (u *RedisSubUsecase) Publish(ctx context.Context, topic string, MessageChan <-chan *pubmessage.PublishMessage) {
	log.Printf("[Pub] start publish %v", topic)
	for {
		message := <-MessageChan
		log.Printf("[Pub] send message: %v", message)
		jsonMessage, err := json.Marshal(message)
		if err != nil {
			log.Printf("pub message convert json error: %v", err)
		}
		u.PubSubRepo.Pub(ctx, topic, []byte(jsonMessage))
	}
}
