package pair

import (
	"context"

	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/ws/client"
)

type PairUsecaseInterface interface {
	SetMessageChan(chan *pubmessage.PublishMessage)
	Add(client *client.Client)
	Run(ctx context.Context)
}
