package message

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/img21326/fb_chat/structure/message"
	"gorm.io/gorm"
)

type MessageRepo struct {
	DB *gorm.DB
}

func NewMessageRepo(db *gorm.DB, redis *redis.Client) MessageRepoInterface {
	return &MessageRepo{
		DB: db,
	}
}

func (r *MessageRepo) Save(ctx context.Context, m *message.Message) {
	r.DB.WithContext(ctx).Create(&m)
}
