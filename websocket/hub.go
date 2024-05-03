package websocket

import (
	"chat_system/websocket/pubsub"
	"context"
	"encoding/json"
	"log"
	"sync"
)

type HubInterface interface {
	AddClient(client ClientInterface)
	RemoveClient(id string)
	FindClient(id string) ClientInterface

	SendMsgToClient(receiver string, msg Msg)
	SendMsgToAll(msg Msg)
	SendMsgToBroker(msg Msg)
	ReceiveMsgFromLocalClient(msg Msg)
	ReceiveMsgFromBroker()

	ConsumeMsg() (msg Msg)
}

type Hub struct {
	clientMap map[string]ClientInterface

	ctx         context.Context
	broker      pubsub.PubSubInterface
	receiveChan chan Msg
	mutex       *sync.Mutex
}

func NewHub(ctx context.Context, opt ...func(*Hub)) HubInterface {
	hub := &Hub{
		clientMap: make(map[string]ClientInterface),

		ctx:         ctx,
		receiveChan: make(chan Msg, 1024),
		mutex:       &sync.Mutex{},
	}

	for _, f := range opt {
		f(hub)
	}

	go hub.ReceiveMsgFromBroker()
	return hub
}

func (h *Hub) SetPubSub(pubSub pubsub.PubSubInterface) {
	h.broker = pubSub
}

func (h *Hub) AddClient(client ClientInterface) {
	h.lockDecorator(func() {
		h.clientMap[client.ID()] = client
	})
}

func (h *Hub) RemoveClient(id string) {
	h.lockDecorator(func() {
		delete(h.clientMap, id)
	})
}

func (h *Hub) FindClient(id string) ClientInterface {
	var client ClientInterface = nil
	h.lockDecorator(func() {
		client = h.clientMap[id]
	})

	return client
}

func (h *Hub) SendMsgToClient(receiver string, msg Msg) {
	msg.Receiver = receiver
	msg.Type = MsgTypeSpecificUser
	client := h.FindClient(receiver)
	if client == nil && msg.From == MsgFromLocal {
		h.SendMsgToBroker(msg)
		return
	}
	client.SendMsg(msg.Msg)
}

func (h *Hub) SendMsgToAll(msg Msg) {
	msg.Type = MsgTypeBroadcast
	for _, client := range h.clientMap {
		if client.ID() == msg.Sender {
			continue
		}
		client.SendMsg(msg.Msg)
	}
	if msg.From == MsgFromLocal {
		h.SendMsgToBroker(msg)
	}
}

func (h *Hub) SendMsgToBroker(msg Msg) {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[SendMsgToBroker] Failed to marshal message: %v\n", err)
		return
	}

	err = h.broker.Publish(msgByte)
	if err != nil {
		log.Printf("[SendMsgToBroker] Failed to publish message: %v\n", err)
	}
}

func (h *Hub) ReceiveMsgFromLocalClient(msg Msg) {
	h.receiveChan <- msg
}

func (h *Hub) ReceiveMsgFromBroker() {
	for {
		select {
		case <-h.ctx.Done():
			return
		default:
			receiveMsg := h.broker.Subscribe()
			msg := Msg{}
			err := json.Unmarshal(receiveMsg, &msg)
			if err != nil {
				log.Printf("[ReceiveMsgFromBroker] Failed to unmarshal message: %v\n", err)
				continue
			}
			if (msg.Receiver == "" || h.FindClient(msg.Receiver) == nil) && msg.Type == MsgTypeSpecificUser {
				log.Printf("[ReceiveMsgFromBroker] Receiver not found: %v\n, msg is: %s", msg.Receiver, string(msg.Msg))
				continue
			}
			msg.From = MsgFromBroker
			if msg.Type == MsgTypeSpecificUser {
				h.SendMsgToClient(msg.Receiver, msg)
			} else {
				h.SendMsgToAll(msg)
			}
		}
	}
}

func (h *Hub) ConsumeMsg() (msg Msg) {
	return <-h.receiveChan
}

func (h *Hub) lockDecorator(f func()) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	f()
}
