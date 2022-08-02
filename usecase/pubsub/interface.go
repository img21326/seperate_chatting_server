package pubsub

import (
	"context"

	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
)

type SubMessageUsecaseInterface interface {
	Subscribe(ctx context.Context, topic string) func() ([]byte, error)
	Publish(ctx context.Context, topic string, Message *pubmessage.PublishMessage)
}
