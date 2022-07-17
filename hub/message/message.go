package message

import (
	"context"
	"log"

	"github.com/img21326/fb_chat/structure/message"
	MessageUsecase "github.com/img21326/fb_chat/usecase/message"
)

type MessageHub struct {
	SaveMessageChan chan *message.Message
	MessageUsecase  MessageUsecase.MessageUsecaseInterface
}

func NewMessageHub(messageUsecase MessageUsecase.MessageUsecaseInterface,
	saveMessageChan chan *message.Message) *MessageHub {
	return &MessageHub{
		MessageUsecase:  messageUsecase,
		SaveMessageChan: saveMessageChan,
	}
}

func (h *MessageHub) Run(ctx context.Context) {
	log.Printf("[MessageUsecase] start\n")
	for {
		select {
		case <-ctx.Done():
			return
		case mes := <-h.SaveMessageChan:
			log.Printf("[MessageUsecase] save message: %+v", mes)
			c := context.Background()
			h.MessageUsecase.Save(c, mes)
		}
	}
}
