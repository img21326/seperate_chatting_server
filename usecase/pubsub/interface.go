package pubsub

import (
	"context"

	repo "github.com/img21326/fb_chat/repo/pubsub"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
)

type SubMessageUsecaseInterface interface {
	Subscribe(ctx context.Context, topic string) repo.SubscribeInterface
	Publish(ctx context.Context, topic string, Message *pubmessage.PublishMessage)
}
