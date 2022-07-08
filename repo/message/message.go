package message

import (
	"context"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type MessageRepo struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewMessageRepo(db *gorm.DB, redis *redis.Client) MessageRepoInterface {
	return &MessageRepo{
		DB:    db,
		Redis: redis,
	}
}

func (r *MessageRepo) Save(ctx context.Context, m *MessageModel) {
	r.DB.WithContext(ctx).Create(&m)
}

func (r *MessageRepo) Publish(ctx context.Context, message []byte) error {
	return r.Redis.Publish(ctx, "message", message).Err()
}
