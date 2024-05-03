package handler

import (
	requestmsg "chat_system/consumer/request_msg"
	"context"

	redis "github.com/redis/go-redis/v9"
)

type LeaveRoomHandler struct {
	baseHandler *BaseHandler
	redis       *redis.Client
}

func NewLeaveRoomHandler(baseHandler *BaseHandler, redis *redis.Client) MessageTypeHandlerInterface {
	return &LeaveRoomHandler{
		baseHandler: baseHandler,
		redis:       redis,
	}
}

func (handler *LeaveRoomHandler) Perform(receiveMsg requestmsg.Msg) error {
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

	err = handler.closeRoom(user)
	if err != nil {
		return err
	}

	msg, err := NewResponseMsg(LeaveRoom, "使用者已離開房間").ToJSON()
	if err != nil {
		return err
	}
	handler.baseHandler.SendServerMsgToUser(user, msg)
	handler.baseHandler.SendServerMsgToUser(pairedUser, msg)

	return nil
}

func (handler *LeaveRoomHandler) closeRoom(user string) error {
	room, err := handler.redis.Get(context.Background(), user).Result()
	if err != nil {
		return err
	}

	_, err = handler.redis.Del(context.Background(), user).Result()
	if err != nil {
		return err
	}

	_, err = handler.redis.Del(context.Background(), room).Result()
	if err != nil {
		return err
	}

	return nil
}
