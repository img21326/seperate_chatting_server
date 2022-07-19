package message

import (
	"context"
	"log"

	"github.com/img21326/fb_chat/structure/message"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	MessageUsecase "github.com/img21326/fb_chat/usecase/message"
	"github.com/img21326/fb_chat/ws/client"
)

type MessageHub struct {
	SaveMessageChan    chan *message.Message
	UnRegisterChan     chan *client.Client
	ReceiveMessageChan chan *pubmessage.PublishMessage

	MessageUsecase MessageUsecase.MessageUsecaseInterface
}

func NewMessageHub(messageUsecase MessageUsecase.MessageUsecaseInterface,
	unRegisterChan chan *client.Client,
	receiveMessageChan chan *pubmessage.PublishMessage) *MessageHub {
	return &MessageHub{
		MessageUsecase:     messageUsecase,
		UnRegisterChan:     unRegisterChan,
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
		case receiveMessage := <-h.ReceiveMessageChan:
			// sendMessage := pubmessage.SendToUserMessage{
			// 	Type:    receiveMessage.Type,
			// 	Payload: receiveMessage.Payload,
			// }
			// jsonMessage, err := json.Marshal(receiveMessage)
			// if err != nil {
			// log.Printf("[WebsocketUsecase] conver receive message err: %v", err)
			// }
			sender, receiver := h.MessageUsecase.GetOnlineClients(receiveMessage.SendFrom, receiveMessage.SendTo)
			// receiveClient, errReceiveClient := u.LocalOnlineRepo.FindUserByID(receiveMessage.SendTo)
			// sendClient, errSendClient := u.LocalOnlineRepo.FindUserByID(receiveMessage.SendFrom)

			if sender == nil && receiver == nil {
				// no client need to hanlde
				continue
			}

			// it will send pair success message to each client
			// so just handler receive side
			if receiveMessage.Type == "pairSuccess" {
				h.MessageUsecase.HandlePairSuccessMessage(receiver, receiveMessage)
				// payload := receiveMessage.Payload.(string)
				// uuid, err := uuid.Parse(payload)
				// if err != nil {
				// 	log.Printf("[WebsocketUsecase] pair success payload convert uuid error: %+v", err)
				// 	continue
				// }
				// if errReceiveClient == nil {
				// 	receiveClient.RoomId = uuid
				// 	receiveClient.PairId = receiveMessage.SendFrom
				// 	receiveClient.Send <- jsonMessage
				// 	log.Printf("[WebsocketUsecase] pair success: %v", receiveClient.User.ID)
				// }
				continue
			}

			// send message to receiver and ack message for sender
			// also sender side need to save message
			if receiveMessage.Type == "pairError" || receiveMessage.Type == "message" || receiveMessage.Type == "leave" {
				h.MessageUsecase.HandleClientOnMessage(sender, receiver, receiveMessage)
				// if receiver != nil {
				// 	receiver.Send <- jsonMessage
				// }
				// // ack message
				// if sender != nil {
				// 	sender.Send <- jsonMessage
				// 	// which server send message, which server should save it.
				// 	if receiveMessage.Type == "message" {
				// 		receiveM, err := json.Marshal(receiveMessage.Payload)
				// 		if err != nil {
				// 			log.Printf("[WebsocketUsecase] save message convert to json err: %v", err)
				// 		}
				// 		var M message.Message
				// 		err = json.Unmarshal(receiveM, &M)
				// 		if err != nil {
				// 			log.Printf("[WebsocketUsecase] save message convert to struct err: %v", err)
				// 		}
				// 		h.SaveMessageChan <- &M
				// 	}
				// }
			}
			if receiveMessage.Type == "leave" {
				h.MessageUsecase.HandleLeaveMessage(sender, receiver, h.UnRegisterChan)
				// if receiver != nil {
				// 	log.Printf("[WebsocketUsecase] send leave message by user %v\n", receiveMessage.SendFrom)
				// 	// 收訊者離開處理
				// 	u.refreshRoomUser(receiveClient)
				// 	u.UnRegisterClientChan <- receiveClient
				// }

				// // 發送者離開處理
				// if errSendClient == nil {
				// 	c := context.Background()
				// 	u.RoomRepo.Close(c, receiveClient.RoomId)
				// 	u.refreshRoomUser(sendClient)
				// 	u.UnRegisterClientChan <- sendClient
				// }
			}
		}
	}
}
