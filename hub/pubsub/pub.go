package pubsub

import (
	"context"
	"log"

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

func (h *PubHub) Run(ctx context.Context, topic string, MessageChan chan *pubmessage.PublishMessage) {
	log.Printf("[Pub] start publish %v", topic)
	for {
		select {
		case <-ctx.Done():
			close(MessageChan)
			log.Printf("[Pub] close channel\n")
			n := len(MessageChan)
			for i := 0; i < n; i++ {
				ctx := context.Background()
				h.PubUsecase.Publish(ctx, topic, <-MessageChan)
			}
			log.Printf("[Pub] finished all queue\n")
			return
		case message := <-MessageChan:
			ctx := context.Background()
			h.PubUsecase.Publish(ctx, topic, message)
		}
	}
}
