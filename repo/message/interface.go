package message

import (
	"context"

	"github.com/img21326/fb_chat/structure/message"
)

type MessageRepoInterface interface {
	Save(context.Context, *message.Message)
}
