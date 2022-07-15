package message

import (
	"context"

	"github.com/img21326/fb_chat/structure/message"
	"github.com/img21326/fb_chat/structure/user"
)

type MessageUsecaseInterface interface {
	SetMessageChan(m chan *message.Message)
	LastByUserID(ctx context.Context, userID uint, c int) (messages []*message.Message, err error) // 從使用者找房間 最後聊天紀錄
	LastByMessageID(ctx context.Context, user *user.User, lastMessageID uint, c int) ([]*message.Message, error)
	Save(*message.Message)
	Run(ctx context.Context)
}
