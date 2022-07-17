package message

import (
	"context"

	"github.com/img21326/fb_chat/structure/message"
)

type MessageUsecaseInterface interface {
	LastByUserID(ctx context.Context, userID uint, c int) (messages []*message.Message, err error) // 從使用者找房間 最後聊天紀錄
	LastByMessageID(ctx context.Context, userID uint, lastMessageID uint, c int) ([]*message.Message, error)
	Save(ctx context.Context, message *message.Message)
}
