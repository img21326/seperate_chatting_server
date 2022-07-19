package message

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/img21326/fb_chat/structure/message"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	MessageUsecase "github.com/img21326/fb_chat/usecase/message"
)

type MessageHub struct {
	SaveMessageChan    chan *message.Message
	MessageUsecase     MessageUsecase.MessageUsecaseInterface
	ReceiveMessageChan chan *pubmessage.PublishMessage
}

func NewMessageHub(messageUsecase MessageUsecase.MessageUsecaseInterface,
	receiveMessageChan chan *pubmessage.PublishMessage) *MessageHub {
	return &MessageHub{
		MessageUsecase:     messageUsecase,
		SaveMessageChan:    make(chan *message.Message, 4096),
		ReceiveMessageChan: receiveMessageChan,
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
			sendClient, errSendClient := u.LocalOnlineRepo.FindUserByID(receiveMessage.SendFrom)

			if errReceiveClient != nil && errSendClient != nil {
				// 這個伺服器皆沒有要處理的使用者
				continue
			}
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
				if errReceiveClient == nil {
					receiveClient.Send <- jsonMessage
				}
				// ack message
				if errSendClient == nil {
					sendClient.Send <- jsonMessage
					if receiveMessage.Type == "message" {
						receiveM, err := json.Marshal(receiveMessage.Payload)
						if err != nil {
							log.Printf("[WebsocketUsecase] save message convert to json err: %v", err)
						}
						var M message.Message
						err = json.Unmarshal(receiveM, &M)
						if err != nil {
							log.Printf("[WebsocketUsecase] save message convert to struct err: %v", err)
						}
						u.SaveMessageChan <- &M
					}
				}
			}
			if receiveMessage.Type == "leave" {
				if errReceiveClient == nil {
					log.Printf("[WebsocketUsecase] send leave message by user %v\n", receiveMessage.SendFrom)
					c := context.Background()
					u.RoomRepo.Close(c, receiveClient.RoomId)
					// 收訊者離開處理
					u.refreshRoomUser(receiveClient)
					u.UnRegisterClientChan <- receiveClient
				}

				// 發送者離開處理
				if errSendClient == nil {
					u.refreshRoomUser(sendClient)
					u.UnRegisterClientChan <- sendClient
				}

			}
		}
	}
}
