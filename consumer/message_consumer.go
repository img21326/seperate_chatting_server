package consumer

import (
	"chat_system/consumer/handler"
	requestmsg "chat_system/consumer/request_msg"
	"chat_system/websocket"
	"errors"
)

type MessageConsumerInterface interface {
	ConsumeMsg(websocket.Msg) error
}

type PairingConsumer struct {
	joinRoomHandler  handler.MessageTypeHandlerInterface
	leaveRoomHandler handler.MessageTypeHandlerInterface
	onMessageHandler handler.MessageTypeHandlerInterface
}

func NewPairingConsumer(joinRoomHandler handler.MessageTypeHandlerInterface,
	leaveRoomHandler handler.MessageTypeHandlerInterface,
	onMessageHandler handler.MessageTypeHandlerInterface,
) MessageConsumerInterface {
	return &PairingConsumer{
		joinRoomHandler:  joinRoomHandler,
		leaveRoomHandler: leaveRoomHandler,
		onMessageHandler: onMessageHandler,
	}
}

var ErrNotMatchEvent = errors.New("not match event")

func (p *PairingConsumer) ConsumeMsg(websocketMsg websocket.Msg) error {
	msg, err := requestmsg.NewRequestMsg(websocketMsg)
	if err != nil {
		return err
	}

	var handler handler.MessageTypeHandlerInterface
	switch msg.MsgEvent {
	case requestmsg.JoinRoom:
		handler = p.joinRoomHandler
	case requestmsg.LeaveRoom:
		handler = p.leaveRoomHandler
	case requestmsg.OnMessage, requestmsg.Typing:
		handler = p.onMessageHandler
	}

	if handler == nil {
		return ErrNotMatchEvent
	}

	return handler.Perform(msg)
}
