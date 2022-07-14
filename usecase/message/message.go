package message

import (
	"context"
	"log"

	RepoMessage "github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/structure/message"
)

type MessageUsecase struct {
	MessageChan chan *message.Message
	MessageRepo RepoMessage.MessageRepoInterface
}

func NewMessageUsecase(messageRepo RepoMessage.MessageRepoInterface) MessageUsecaseInterface {
	return &MessageUsecase{
		MessageRepo: messageRepo,
	}
}

func (u *MessageUsecase) Last(ctx context.Context, lastMessageID uint, c int) (messages []*message.Message, err error) {
	lastMessage, err := u.MessageRepo.GetByID(ctx, lastMessageID)
	if err != nil {
		return nil, err
	}
	messages, err = u.MessageRepo.LastsByTime(ctx, lastMessage.Time, c)
	return
}

func (u *MessageUsecase) Save(m *message.Message) {
	u.MessageChan <- m
}

func (u *MessageUsecase) SetMessageChan(m chan *message.Message) {
	u.MessageChan = m
}

func (u *MessageUsecase) Run(ctx context.Context) {
	log.Printf("[MessageUsecase] start\n")
	for {
		select {
		case <-ctx.Done():
			return
		case mes := <-u.MessageChan:
			log.Printf("[MessageUsecase] save message: %+v", mes)
			u.MessageRepo.Save(ctx, mes)
		}
	}
}
