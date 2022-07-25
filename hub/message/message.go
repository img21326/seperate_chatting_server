package message

import (
	"context"
	"log"

	"github.com/img21326/fb_chat/structure/message"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	MessageUsecase "github.com/img21326/fb_chat/usecase/message"
	WSUsecase "github.com/img21326/fb_chat/usecase/ws"
	"github.com/img21326/fb_chat/ws/client"
)

type MessageHub struct {
	SaveMessageChan    chan *message.Message
	ReceiveMessageChan chan *pubmessage.PublishMessage

	MessageUsecase MessageUsecase.MessageUsecaseInterface
	WSUsecase      WSUsecase.WebsocketUsecaseInterface
}

func NewMessageHub(messageUsecase MessageUsecase.MessageUsecaseInterface,
	wsUsecase WSUsecase.WebsocketUsecaseInterface,
	unRegisterChan chan *client.Client,
	receiveMessageChan chan *pubmessage.PublishMessage) *MessageHub {
	return &MessageHub{
		MessageUsecase:     messageUsecase,
		WSUsecase:          wsUsecase,
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
			sender, receiver := h.MessageUsecase.GetOnlineClients(receiveMessage.SendFrom, receiveMessage.SendTo)
			log.Printf("[MessageHub] get user %v, %v", sender, receiver)
			if sender == nil && receiver == nil {
				// no client need to hanlde
				log.Printf("[MessageHub] no user")
				continue
			}

			// it will send pair success message to each client
			// so just handler receive side
			if receiveMessage.Type == "pairSuccess" {
				err := h.MessageUsecase.HandlePairSuccessMessage(receiver, receiveMessage)
				if err != nil {
					log.Printf("[MessageHub] HandlePairSuccessMessage err: %v", err)
				}
				continue
			}

			// send message to receiver and ack message for sender
			// also sender side need to save message
			if receiveMessage.Type == "pairError" || receiveMessage.Type == "message" || receiveMessage.Type == "leave" {
				err := h.MessageUsecase.HandleClientOnMessage(sender, receiver, receiveMessage, h.SaveMessageChan)
				if err != nil {
					log.Printf("[MessageHub] HandleClientOnMessage err: %v", err)
				}
			}

			if receiveMessage.Type == "leave" {
				h.MessageUsecase.HandleLeaveMessage(sender, receiver, h.WSUsecase.UnRegister)
			}
		}
	}
}
