package message

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	localOnline "github.com/img21326/fb_chat/repo/local_online"
	RepoMessage "github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/repo/room"
	errorStruct "github.com/img21326/fb_chat/structure/error"
	"github.com/img21326/fb_chat/structure/message"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"

	"github.com/img21326/fb_chat/ws/client"
)

type MessageUsecase struct {
	MessageChan     chan *message.Message
	MessageRepo     RepoMessage.MessageRepoInterface
	RoomRepo        room.RoomRepoInterface
	LocalOnlineRepo localOnline.LocalOnlineRepoInterface
}

func NewMessageUsecase(messageRepo RepoMessage.MessageRepoInterface,
	roomRepo room.RoomRepoInterface, localOnlineRepo localOnline.LocalOnlineRepoInterface) MessageUsecaseInterface {
	return &MessageUsecase{
		MessageRepo:     messageRepo,
		RoomRepo:        roomRepo,
		LocalOnlineRepo: localOnlineRepo,
	}
}

func (u *MessageUsecase) LastByUserID(ctx context.Context, userID uint, c int) (messages []*message.Message, err error) {
	room, err := u.RoomRepo.FindByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	if room.Close {
		return nil, errorStruct.RoomIsClose
	}
	messages, err = u.MessageRepo.LastsByRoomID(ctx, room.UUID, c)
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
		return nil, errorStruct.UserNotInThisRoom
	}
	if room.Close {
		return nil, errorStruct.RoomIsClose
	}
	messages, err = u.MessageRepo.LastsByTime(ctx, lastMessage.RoomId, lastMessage.Time, c)
	return
}

func (u *MessageUsecase) Save(ctx context.Context, message *message.Message) {
	u.MessageRepo.Save(ctx, message)
}

func (u *MessageUsecase) GetOnlineClients(senderID uint, receiverID uint) (sender *client.Client, receiver *client.Client) {
	sender, _ = u.LocalOnlineRepo.FindUserByID(senderID)
	receiver, _ = u.LocalOnlineRepo.FindUserByID(receiverID)
	return
}

func (u *MessageUsecase) HandlePairSuccessMessage(receiver *client.Client, receiveMessage *pubmessage.PublishMessage) error {
	if receiver == nil {
		return errorStruct.ClientNotInHost
	}
	payload := receiveMessage.Payload.(string)
	uuid, err := uuid.Parse(payload)
	if err != nil {
		log.Printf("[MessageUsecase] HandlePairSuccessMessage convert uuid err: %v", err)
		return err
	}
	jsonMessage, err := json.Marshal(receiveMessage)
	if err != nil {
		log.Printf("[MessageUsecase] HandlePairSuccessMessage convert receive message err: %v", err)
		return err
	}
	receiver.RoomId = uuid
	receiver.PairId = receiveMessage.SendFrom
	err = receiver.Send.Push(jsonMessage)
	if err != nil {
		log.Printf("[MessageUsecase] HandlePairSuccessMessage receiver send message err: %v", err)
	}
	return nil
}

func (u *MessageUsecase) HandleClientOnMessage(sender *client.Client, receiver *client.Client, receiveMessage *pubmessage.PublishMessage, saveMessageChan chan *message.Message) error {
	jsonMessage, err := json.Marshal(receiveMessage)
	if err != nil {
		log.Printf("[MessageUsecase] HandleClientOnMessage convert receive message err: %v", err)
		return err
	}
	if receiver != nil {
		err := receiver.Send.Push(jsonMessage)
		if err != nil {
			log.Printf("[MessageUsecase] HandleClientOnMessage receiver send message err: %v", err)
		}
	}
	// ack message
	if sender != nil {
		err := sender.Send.Push(jsonMessage)
		if err != nil {
			log.Printf("[MessageUsecase] HandleClientOnMessage sender send message err: %v", err)
		}
		// which server send message, which server should save it.
		if receiveMessage.Type == "message" {
			receiveM, err := json.Marshal(receiveMessage.Payload)
			if err != nil {
				log.Printf("[MessageUsecase] HandleClientOnMessage save message convert to json err: %v", err)
				return err
			}
			var M message.Message
			err = json.Unmarshal(receiveM, &M)
			if err != nil {
				log.Printf("[MessageUsecase] HandleClientOnMessage save message convert to struct err: %v", err)
				return err
			}
			saveMessageChan <- &M
		}
	}
	return nil
}

func (u *MessageUsecase) refreshRoomUser(client *client.Client) {
	client.RoomId = uuid.Nil
	client.PairId = 0
}

func (u *MessageUsecase) HandleLeaveMessage(sender *client.Client, receiver *client.Client, unRegisterFunc func(ctx context.Context, client *client.Client)) error {
	c := context.Background()

	// 發送者離開處理
	if sender != nil {
		u.RoomRepo.Close(c, sender.RoomId)
		u.refreshRoomUser(sender)
		unRegisterFunc(c, sender)
	}
	// 收訊者離開處理
	if receiver != nil {
		u.refreshRoomUser(receiver)
		unRegisterFunc(c, receiver)
	}
	return nil
}
