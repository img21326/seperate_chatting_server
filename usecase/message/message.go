package message

import (
	"context"
	"errors"
	"log"

	RepoMessage "github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/repo/room"
	"github.com/img21326/fb_chat/structure/message"
	"github.com/img21326/fb_chat/structure/user"
)

type MessageUsecase struct {
	MessageChan chan *message.Message
	MessageRepo RepoMessage.MessageRepoInterface
	RoomRepo    room.RoomRepoInterface
}

func NewMessageUsecase(messageRepo RepoMessage.MessageRepoInterface, roomRepo room.RoomRepoInterface) MessageUsecaseInterface {
	return &MessageUsecase{
		MessageRepo: messageRepo,
		RoomRepo:    roomRepo,
	}
}

func (u *MessageUsecase) LastByUserID(ctx context.Context, userID uint, c int) (messages []*message.Message, err error) {
	room, err := u.RoomRepo.FindByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	messages, err = u.MessageRepo.LastsByRoomID(ctx, room.ID, 30)
	return
}

func (u *MessageUsecase) LastByMessageID(ctx context.Context, user *user.User, lastMessageID uint, c int) (messages []*message.Message, err error) {
	lastMessage, err := u.MessageRepo.GetByID(ctx, lastMessageID)
	if err != nil {
		return nil, err
	}
	room, err := u.RoomRepo.FindByUserId(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if room.UserId1 != user.ID && room.UserId2 != user.ID {
		return nil, errors.New("UserNotInThisRoom")
	}
	messages, err = u.MessageRepo.LastsByTime(ctx, lastMessage.RoomId, lastMessage.Time, c)
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
