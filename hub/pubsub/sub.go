package pubsub

import (
	"context"

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

func (h *SubHub) Run(ctx context.Context, topic string, MessageChan chan<- *pubmessage.PublishMessage) {
	h.PubUsecase.Subscribe(ctx, topic, func(pm *pubmessage.PublishMessage) {
		MessageChan <- pm
	})
}
