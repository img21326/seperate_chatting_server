package handler

import (
	requestmsg "chat_system/consumer/request_msg"
	"chat_system/websocket"
)

type OnMessageHandler struct {
	baseHandler *BaseHandler
	hub         websocket.HubInterface
}

func NewOnMessageHandler(baseHandler *BaseHandler, hub websocket.HubInterface) MessageTypeHandlerInterface {
	return &OnMessageHandler{
		baseHandler: baseHandler,
		hub:         hub,
	}
}

func (handler *OnMessageHandler) Perform(receiveMsg requestmsg.Msg) error {
	user := receiveMsg.Sender

	if !handler.baseHandler.isInRoom(user) {
		msg, err := NewResponseMsg(NotInRoom, "您不在房間裡").ToJSON()
		if err != nil {
			return err
		}
		handler.baseHandler.SendServerMsgToUser(user, msg)
		return nil
	}

	pairedUser, err := handler.baseHandler.getPairedUser(user)
	if err != nil {
		return err
	}
	msg, err := NewResponseMsg(OnMessage, string(receiveMsg.Msg)).ToJSON()
	if err != nil {
		return err
	}
	handler.hub.SendMsgToClient(pairedUser, websocket.Msg{
		Sender: user,
		Msg:    msg,
	})
	return nil
}
