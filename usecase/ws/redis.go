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

func (u *RedisWebsocketUsecase) refreshRoomUser(client *client.Client) {
	client.ContinueLoop = false
	client.RoomId = uuid.Nil
	client.PairId = 0
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
			log.Printf("[WebsocketUsecase] Register user: %v \n", client.User.ID)
		case client := <-u.UnRegisterClientChan:
			u.LocalOnlineRepo.UnRegister(client)
			u.OnlineRepo.UnRegister(ctx, client.User.ID)
			_ = client.Conn.Close()
			close(client.Send)
			log.Printf("[WebsocketUsecase] unRegister user: %v \n", client.User.ID)
		case receiveMessage := <-u.ReceiveMessageChan:
			// sendMessage := pubmessage.SendToUserMessage{
			// 	Type:    receiveMessage.Type,
			// 	Payload: receiveMessage.Payload,
			// }
			jsonMessage, err := json.Marshal(receiveMessage)
			if err != nil {
				log.Printf("[WebsocketUsecase] conver receive message err: %v", err)
			}
			receiveClient, errReceiveClient := u.LocalOnlineRepo.FindUserByID(receiveMessage.SendTo)

			sendClient, errPairClient := u.LocalOnlineRepo.FindUserByID(receiveMessage.SendFrom)
			// pairSuccess會發給兩個使用者 所以不用care
			if receiveMessage.Type == "pairSuccess" {
				payload := receiveMessage.Payload.(string)
				uuid, err := uuid.Parse(payload)
				if err != nil {
					log.Printf("[WebsocketUsecase] pair success payload convert uuid error: %+v", err)
					continue
				}
				if errReceiveClient == nil {
					receiveClient.RoomId = uuid
					receiveClient.PairId = receiveMessage.SendFrom
					receiveClient.Send <- jsonMessage
					log.Printf("[WebsocketUsecase] pair success: %v", receiveClient.User.ID)
				}
				continue
			}
			if receiveMessage.Type == "pairError" || receiveMessage.Type == "message" || receiveMessage.Type == "leave" {
				// ask message
				if errReceiveClient == nil {
					receiveClient.Send <- jsonMessage
				}
				if errPairClient == nil {
					sendClient.Send <- jsonMessage
				}
			}
			if receiveMessage.Type == "leave" {
				if errReceiveClient == nil {
					log.Printf("[WebsocketUsecase] send leave message by user %v\n", receiveMessage.SendFrom)
					u.RoomRepo.Close(receiveClient.RoomId)
					// 收訊者斷線處理
					u.refreshRoomUser(receiveClient)
					u.UnRegisterClientChan <- receiveClient
				}

				// 發送者斷線處理
				if errPairClient == nil {
					u.refreshRoomUser(sendClient)
					u.UnRegisterClientChan <- sendClient
				}

			}
		}
	}
}
