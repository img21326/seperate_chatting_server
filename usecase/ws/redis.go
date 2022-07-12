package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	localonline "github.com/img21326/fb_chat/repo/local_online"
	"github.com/img21326/fb_chat/repo/online"
	RepoRoom "github.com/img21326/fb_chat/repo/room"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/structure/room"
	"github.com/img21326/fb_chat/ws/client"
)

type RedisWebsocketUsecase struct {
	RegisterClientChan   chan *client.Client
	UnRegisterClientChan chan *client.Client
	ReceiveMessageChan   chan *pubmessage.PublishMessage
	LocalOnlineRepo      localonline.OnlineRepoInterface
	OnlineRepo           online.OnlineRepoInterface
	RoomRepo             RepoRoom.RoomRepoInterface
}

func NewRedisWebsocketUsecase(
	// registerClientChan <-chan *client.Client, unregisterCLientChan chan *client.Client,
	// receiveMessage <-chan *pubmessage.PublishMessage,
	localOnlineRepo localonline.OnlineRepoInterface, onlineRepo online.OnlineRepoInterface,
	roomRepo RepoRoom.RoomRepoInterface,
) WebsocketUsecaseInterface {
	return &RedisWebsocketUsecase{
		// RegisterClientChan:   registerClientChan,
		// UnRegisterClientChan: unregisterCLientChan,
		// ReceiveMessageChan:   receiveMessage,
		RegisterClientChan:   make(chan *client.Client, 1024),
		UnRegisterClientChan: make(chan *client.Client, 1024),
		ReceiveMessageChan:   make(chan *pubmessage.PublishMessage, 1024),
		LocalOnlineRepo:      localOnlineRepo,
		OnlineRepo:           onlineRepo,
		RoomRepo:             roomRepo,
	}
}

func (u *RedisWebsocketUsecase) FindRoomByUserId(ctx context.Context, userID uint) (*room.Room, error) {
	return u.RoomRepo.FindByUserId(ctx, userID)
}

func (u *RedisWebsocketUsecase) ReceiveMessage(message *pubmessage.PublishMessage) {
	u.ReceiveMessageChan <- message
}

func (u *RedisWebsocketUsecase) UnRegister(client *client.Client) {
	u.UnRegisterClientChan <- client
}

func (u *RedisWebsocketUsecase) Register(client *client.Client) {
	u.RegisterClientChan <- client
}

func (u *RedisWebsocketUsecase) Run(ctx context.Context) {
	log.Printf("[WebsocketUsecase] start")
	for {
		select {
		case client := <-u.RegisterClientChan:
			u.LocalOnlineRepo.Register(client)
			u.OnlineRepo.Register(ctx, client.User.ID)
		case client := <-u.UnRegisterClientChan:
			u.LocalOnlineRepo.UnRegister(client)
			u.OnlineRepo.UnRegister(ctx, client.User.ID)
			client.Conn.Close()
			close(client.Send)
		case receiveMessage := <-u.ReceiveMessageChan:
			sendMessage := pubmessage.SendToUserMessage{
				Type:    receiveMessage.Type,
				Payload: receiveMessage.Payload,
			}
			jsonMessage, err := json.Marshal(sendMessage)
			if err != nil {
				log.Printf("[WebsocketUsecase] conver receive message err: %v", err)
			}
			client, err := u.LocalOnlineRepo.FindUserByFbID(receiveMessage.SendTo)
			if err != nil {
				log.Printf("[WebsocketUsecase] receive message not found online user")
				continue
			}
			if receiveMessage.Type == "pairSuccess" {
				payload := receiveMessage.Payload.(string)
				uuid, err := uuid.Parse(payload)
				if err != nil {
					log.Printf("[WebsocketUsecase] pair success payload convert uuid error: %+v", err)
					continue
				}
				client.RoomId = uuid
				client.PairId = receiveMessage.SendFrom
				client.Send <- jsonMessage
				log.Printf("[WebsocketUsecase] pair success: %v", client.User.ID)
				continue
			}
			if receiveMessage.Type == "pairError" || receiveMessage.Type == "message" || receiveMessage.Type == "leave" {
				client.Send <- jsonMessage
			}
			if receiveMessage.Type == "leave" {
				log.Printf("[WebsocketUsecase] send leave message by user %v\n", client.User.ID)
				u.RoomRepo.Close(client.RoomId)
				client.RoomId = uuid.Nil
				client.PairId = 0
				u.UnRegisterClientChan <- client
			}
		}
	}
}
