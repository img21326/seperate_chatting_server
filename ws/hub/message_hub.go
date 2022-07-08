package hub

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/usecase/hub"
)

type MessageQueue struct {
	SendMessage chan *message.MessageModel
	Close       chan uuid.UUID
	HubUsecase  hub.HubUsecaseInterface
}

func (q *MessageQueue) Run() {
	log.Printf("[MessageQueue] start \n")
	for {
		select {
		case message := <-q.SendMessage:
			log.Printf("[MessageQueue] new message: %+v \n", message)
			ctx := context.Background()
			q.HubUsecase.SaveMesssage(ctx, message)
			log.Printf("[MessageQueue] save message done \n")
		}
	}
}
