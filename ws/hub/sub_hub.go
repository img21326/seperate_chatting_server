package hub

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/img21326/fb_chat/ws/messageType"
)

type SubHub struct {
	ReceiveChan chan messageType.PublishMessage
	Redis       *redis.Client
}

func (h *SubHub) MessageController(ctx context.Context) {
	log.Printf("[subMessageController] start")
	subscriber := h.Redis.Subscribe(ctx, "message")
	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("sub message receive error: %v", err)
		}
		var redisMessage messageType.PublishMessage

		if err := json.Unmarshal([]byte(msg.Payload), &redisMessage); err != nil {
			log.Printf("pubsub message json load error: %v", err)
		}
		h.ReceiveChan <- redisMessage
	}
}
