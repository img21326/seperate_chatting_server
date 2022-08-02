package message

import (
	"context"
	"log"
	"time"

	"github.com/img21326/fb_chat/structure/message"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	MessageUsecase "github.com/img21326/fb_chat/usecase/message"
	WSUsecase "github.com/img21326/fb_chat/usecase/ws"
)

type MessageHub struct {
	SaveMessageChan    chan *message.Message
	ReceiveMessageChan chan *pubmessage.PublishMessage
	PubMessageChan     chan *pubmessage.PublishMessage

	MessageUsecase MessageUsecase.MessageUsecaseInterface
	WSUsecase      WSUsecase.WebsocketUsecaseInterface
}

func NewMessageHub(messageUsecase MessageUsecase.MessageUsecaseInterface,
	wsUsecase WSUsecase.WebsocketUsecaseInterface,
	receiveMessageChan chan *pubmessage.PublishMessage) *MessageHub {
	return &MessageHub{
		MessageUsecase:     messageUsecase,
		WSUsecase:          wsUsecase,
		SaveMessageChan:    make(chan *message.Message, 4096),
		ReceiveMessageChan: receiveMessageChan,
	}
}

func (h *MessageHub) SaveMessage(mes *message.Message) {
	log.Printf("[MessageHub] save message: %+v", mes)
	c := context.Background()
	h.MessageUsecase.Save(c, mes)
}

func (h *MessageHub) HandleMessage(receiveMessage *pubmessage.PublishMessage) {
	// sendMessage := pubmessage.SendToUserMessage{
	// 	Type:    receiveMessage.Type,
	// 	Payload: receiveMessage.Payload,
	// }
	sender, receiver := h.MessageUsecase.GetOnlineClients(receiveMessage.SendFrom, receiveMessage.SendTo)
	log.Printf("[MessageHub] get user %v, %v", sender, receiver)
	if sender == nil && receiver == nil {
		// no client need to hanlde
		log.Printf("[MessageHub] no user")
		return
	}

	// it will send pair success message to each client
	// so just handler receive side
	if receiveMessage.Type == "pairSuccess" {
		err := h.MessageUsecase.HandlePairSuccessMessage(receiver, receiveMessage)
		if err != nil {
			log.Printf("[MessageHub] HandlePairSuccessMessage err: %v", err)
		}
		return
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
		err := h.MessageUsecase.HandleLeaveMessage(sender, receiver, h.WSUsecase.UnRegister)
		if err != nil {
			log.Printf("[MessageHub] HandleLeaveMessage err: %v", err)
		}
	}
}

func (h *MessageHub) Run(ctx context.Context, timeOut time.Duration) {
	log.Printf("[MessageHub] start\n")
	for {
		select {
		case <-ctx.Done():
			time.Sleep(time.Duration(timeOut))
			close(h.SaveMessageChan)
			close(h.ReceiveMessageChan)
			log.Printf("[MessageHub] close channel\n")
			// 關閉server後不再處理相關訊息 只儲存已發送(ask)訊息
			n := len(h.SaveMessageChan)
			for i := 0; i < n; i++ {
				h.SaveMessage(<-h.SaveMessageChan)
			}
			// n = len(h.ReceiveMessageChan)
			// for i := 0; i < n; i++ {
			// 	h.HandleMessage(<-h.ReceiveMessageChan)
			// }
			log.Printf("[MessageHub] finished all queue\n")
			return
		case mes := <-h.SaveMessageChan:
			h.SaveMessage(mes)
		case receiveMessage := <-h.ReceiveMessageChan:
			h.HandleMessage(receiveMessage)
		}
	}
}
