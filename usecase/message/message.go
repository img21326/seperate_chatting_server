package message

import (
	"context"
	"errors"

	RepoMessage "github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/repo/room"
	"github.com/img21326/fb_chat/structure/message"
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
	messages, err = u.MessageRepo.LastsByRoomID(ctx, room.ID, c)
	return
}

func (u *MessageUsecase) LastByMessageID(ctx context.Context, userID uint, lastMessageID uint, c int) (messages []*message.Message, err error) {
	lastMessage, err := u.MessageRepo.GetByID(ctx, lastMessageID)
	if err != nil {
		return nil, err
	}
	room, err := u.RoomRepo.FindByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	if room.UserId1 != userID && room.UserId2 != userID {
		return nil, errors.New("UserNotInThisRoom")
	}
	messages, err = u.MessageRepo.LastsByTime(ctx, lastMessage.RoomId, lastMessage.Time, c)
	return
}

func (u *MessageUsecase) Save(ctx context.Context, message *message.Message) {
	u.MessageRepo.Save(ctx, message)
}
