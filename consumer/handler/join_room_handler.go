package handler

import (
	requestmsg "chat_system/consumer/request_msg"
	"context"

	redis "github.com/redis/go-redis/v9"
)

type JoinRoomHandler struct {
	baseHandler *BaseHandler
	redis       *redis.Client
}

func NewJoinRoomHandler(baseHandler *BaseHandler, redis *redis.Client) MessageTypeHandlerInterface {
	return &JoinRoomHandler{
		baseHandler: baseHandler,
		redis:       redis,
	}
}

func (handler *JoinRoomHandler) Perform(receiveMsg requestmsg.Msg) error {
	user := receiveMsg.Sender

	if handler.baseHandler.isInRoom(user) {
		msg, err := NewResponseMsg(InRoom, "您已回到原本的房間").ToJSON()
		if err != nil {
			return err
		}
		handler.baseHandler.SendServerMsgToUser(user, msg)
		return nil
	}

	if handler.userIsWaiting(user) {
		msg, err := NewResponseMsg(WaitingForPairing, "您已在等待配對中").ToJSON()
		if err != nil {
			return err
		}
		handler.baseHandler.SendServerMsgToUser(user, msg)
		return nil
	}

	if handler.hasWaitingUser() {
		waitingUser, err := handler.getWaitingUser()
		if err != nil {
			return err
		}
		handler.createRoom(user, waitingUser)

		msg, err := NewResponseMsg(InRoom, "配對成功，已為您開啟房間").ToJSON()
		if err != nil {
			return err
		}
		handler.baseHandler.SendServerMsgToUser(user, msg)
		handler.baseHandler.SendServerMsgToUser(waitingUser, msg)
		return nil
	}

	msg, err := NewResponseMsg(WaitingForPairing, "等待配對中...").ToJSON()
	if err != nil {
		return err
	}
	err = handler.setWaitingUser(user)
	if err != nil {
		return err
	}
	handler.baseHandler.SendServerMsgToUser(user, msg)
	return nil
}

func (handler *JoinRoomHandler) userIsWaiting(user string) bool {
	position, err := handler.redis.LPos(context.Background(), "waiting", user, redis.LPosArgs{}).Result()
	return err == nil && position != -1
}

func (handler *JoinRoomHandler) hasWaitingUser() bool {
	length, err := handler.redis.LLen(context.Background(), "waiting").Result()
	if err != nil {
		return false
	}
	return length > 0
}

func (handler *JoinRoomHandler) getWaitingUser() (string, error) {
	user, err := handler.redis.LPop(context.Background(), "waiting").Result()
	if err != nil {
		return "", err
	}
	return user, nil
}

func (handler *JoinRoomHandler) createRoom(user1, user2 string) error {
	err := handler.redis.Set(context.Background(), user1, user2, 0).Err()
	if err != nil {
		return err
	}
	return handler.redis.Set(context.Background(), user2, user1, 0).Err()
}

func (handler *JoinRoomHandler) setWaitingUser(user string) error {
	return handler.redis.RPush(context.Background(), "waiting", user).Err()
}
