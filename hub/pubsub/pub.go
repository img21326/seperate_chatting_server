package pubsub

import (
	"context"

	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/usecase/pubsub"
)

type PubHub struct {
	PubUsecase pubsub.SubMessageUsecaseInterface
}

func NewPubHub(pubUsecase pubsub.SubMessageUsecaseInterface) *PubHub {
	return &PubHub{
		PubUsecase: pubUsecase,
	}
}

func (h *PubHub) Run(ctx context.Context, topic string, MessageChan <-chan *pubmessage.PublishMessage) {
	h.PubUsecase.Publish(ctx, topic, MessageChan)
}
