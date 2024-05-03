package handler

import (
	"chat_system/websocket"
	"context"

	redis "github.com/redis/go-redis/v9"
)

type BaseHandler struct {
	hub   websocket.HubInterface
	redis *redis.Client
}

func NewBaseHandler(hub websocket.HubInterface, redis *redis.Client) *BaseHandler {
	return &BaseHandler{
		hub:   hub,
		redis: redis,
	}
}

func (handler *BaseHandler) SendServerMsgToUser(user string, msg []byte) {
	handler.hub.SendMsgToClient(user, websocket.Msg{
		Sender: "SERVER",
		Msg:    msg,
	})
}

func (handler *BaseHandler) isInRoom(user string) bool {
	res, err := handler.redis.Get(context.Background(), user).Result()
	if err != nil {
		return false
	}
	return res != ""
}

func (handler *BaseHandler) getPairedUser(user string) (string, error) {
	pairedUser, err := handler.redis.Get(context.Background(), user).Result()
	if err != nil {
		return "", err
	}
	return pairedUser, nil
}
