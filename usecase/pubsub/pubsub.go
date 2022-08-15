package pubsub

import (
	"context"
	"encoding/json"
	"log"

	repo "github.com/img21326/fb_chat/repo/pubsub"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
)

type SubUsecase struct {
	PubSubRepo repo.PubSubRepoInterface
}

func NewSubUsecase(pubSubRepo repo.PubSubRepoInterface) SubMessageUsecaseInterface {
	return &SubUsecase{
		PubSubRepo: pubSubRepo,
	}
}

func (u *SubUsecase) Subscribe(ctx context.Context, topic string) func() ([]byte, error) {
	subscriber := u.PubSubRepo.Sub(ctx, topic)
	return subscriber
}

func (u *SubUsecase) Publish(ctx context.Context, topic string, message *pubmessage.PublishMessage) {
	log.Printf("[PublishUsecase] send message: %v", message)
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("pub message convert json error: %v", err)
	}
	u.PubSubRepo.Pub(ctx, topic, []byte(jsonMessage))
}
