package pubsub

import (
	"context"
	"encoding/json"
	"log"

	errorStruct "github.com/img21326/fb_chat/structure/error"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"

	"github.com/img21326/fb_chat/usecase/pubsub"
)

type SubHub struct {
	PubUsecase pubsub.SubMessageUsecaseInterface
}

func NewSubHub(pubUsecase pubsub.SubMessageUsecaseInterface) *SubHub {
	return &SubHub{
		PubUsecase: pubUsecase,
	}
}

func (h *SubHub) Run(ctx context.Context, topic string, ReceiveMessageChan chan *pubmessage.PublishMessage) {
	log.Printf("[Sub] start subscribe %v", topic)
	subscriber := h.PubUsecase.Subscribe(ctx, topic)
	// subscriber will be close in repo package
	for {
		msg, err := subscriber()
		log.Printf("[Sub] get message: %v", msg)
		if err != nil {
			log.Printf("[Sub] sub message receive error: %v", err)
			if err == errorStruct.ChannelClosed {
				return
			}
		}
		var Message pubmessage.PublishMessage

		if err := json.Unmarshal(msg, &Message); err != nil {
			log.Printf("pubsub message json load error: %v", err)
		}
		ReceiveMessageChan <- &Message
	}
}
