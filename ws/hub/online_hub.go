package hub

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/usecase/hub"
	"github.com/img21326/fb_chat/ws/client"
	"github.com/img21326/fb_chat/ws/messageType"
)

type OnlineHub struct {
	Register     chan *client.Client
	Unregister   chan *client.Client
	ReceiveChan  chan messageType.PublishMessage
	PublishChan  chan messageType.PublishMessage
	MessageQueue *MessageQueue
	HubUsecase   hub.HubUsecaseInterface
}

func (h *OnlineHub) Run() {
	log.Printf("[onlineHub] start")
	for {
		select {
		case client := <-h.Register:
			h.HubUsecase.RegisterOnline(client)
			log.Printf("[onlineHub] %v register online success\n", client.User.ID)
		case client := <-h.Unregister:
			h.HubUsecase.UnRegisterOnline(client)
			log.Printf("[onlineHub] %v unregister success\n", client.User.ID)
		case receiveMessage := <-h.ReceiveChan:
			log.Printf("[onlineHub] receive message: %+v\n", receiveMessage)
			sendMessage := messageType.SendToUserMessage{
				Type:    receiveMessage.Type,
				Payload: receiveMessage.Payload,
			}
			jsonMessage, _ := json.Marshal(sendMessage)
			user, err := h.HubUsecase.FindOnlineUserByUserID(receiveMessage.SendTo)
			if err != nil {
				log.Printf("[onlineHub] receive message not found online user")
				continue
			}
			if receiveMessage.Type == "pairSuccess" {
				payload := receiveMessage.Payload.(string)
				uuid, err := uuid.Parse(payload)
				if err != nil {
					log.Printf("[onlineHub] pair success payload error: %+v", err)
					continue
				}
				user.RoomId = uuid
				user.PairId = receiveMessage.SendFrom
				user.Send <- jsonMessage
				log.Printf("[onlineHub] pair success: %+v", user)
				continue
			}
			if receiveMessage.Type == "pairError" || receiveMessage.Type == "message" || receiveMessage.Type == "leave" {
				user.Send <- jsonMessage
			}
			if receiveMessage.Type == "leave" {
				log.Printf("[onlineHub] send leave message by user %v\n", user.User.ID)
				h.HubUsecase.CloseRoom(user.RoomId)
				user.RoomId = uuid.Nil
				user.PairId = 0
				user.Conn.Close()
				h.Unregister <- user
			}
		case publishMessage := <-h.PublishChan:
			log.Printf("[onlineHub] publish message: %+v\n", publishMessage)
			if publishMessage.Type == "leave" {
				// close room
				roomID := publishMessage.Payload.(uuid.UUID)
				h.MessageQueue.Close <- roomID
			}
			if publishMessage.Type == "message" {
				// save message
				m := publishMessage.Payload.(message.MessageModel)
				h.MessageQueue.SendMessage <- &m
				// change payload (dont send all of message model)
				type M struct {
					Message string
					Time    time.Time
				}
				mes := M{
					Message: m.Message,
					Time:    m.Time,
				}
				publishMessage.Payload = mes
			}
			ctx := context.Background()
			err := h.HubUsecase.SendMessage(ctx, publishMessage)
			if err != nil {
				log.Printf("[onlineHub]publish message error: %v", err)
				continue
			}
		}
	}
}
