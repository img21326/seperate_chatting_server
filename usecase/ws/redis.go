package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	localonline "github.com/img21326/fb_chat/repo/local_online"
	"github.com/img21326/fb_chat/repo/online"
	"github.com/img21326/fb_chat/repo/room"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/ws/client"
)

type RedisWebsocketUsecase struct {
	RegisterClientChan   <-chan *client.Client
	UnRegisterClientChan chan *client.Client
	ReceiveMessageChan   <-chan *pubmessage.PublishMessage
	LocalOnlineRepo      localonline.OnlineRepoInterface
	OnlineRepo           online.OnlineRepoInterface
	RoomRepo             room.RoomRepoInterface
}

func NewRedisWebsocketUsecase(registerClientChan <-chan *client.Client, unregisterCLientChan chan *client.Client,
	receiveMessage <-chan *pubmessage.PublishMessage,
	localOnlineRepo localonline.OnlineRepoInterface, onlineRepo online.OnlineRepoInterface,
	roomRepo room.RoomRepoInterface,
) WebsocketUsecaseInterface {
	return &RedisWebsocketUsecase{
		RegisterClientChan:   registerClientChan,
		UnRegisterClientChan: unregisterCLientChan,
		ReceiveMessageChan:   receiveMessage,
		LocalOnlineRepo:      localOnlineRepo,
		OnlineRepo:           onlineRepo,
		RoomRepo:             roomRepo,
	}
}

func (u *RedisWebsocketUsecase) Run(ctx context.Context) {
	for {
		select {
		case client := <-u.RegisterClientChan:
			u.LocalOnlineRepo.Register(client)
			u.OnlineRepo.Register(ctx, client.User.ID)
		case client := <-u.UnRegisterClientChan:
			u.LocalOnlineRepo.UnRegister(client)
			u.OnlineRepo.UnRegister(ctx, client.User.ID)
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
				log.Printf("[onlineHub] send leave message by user %v\n", client.User.ID)
				u.RoomRepo.Close(client.RoomId)
				client.RoomId = uuid.Nil
				client.PairId = 0
				client.Conn.Close()
				u.UnRegisterClientChan <- client
			}
		}
	}
}
