package message

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/img21326/fb_chat/structure/message"
	"gorm.io/gorm"
)

type MessageRepo struct {
	DB *gorm.DB
}

func NewMessageRepo(db *gorm.DB) MessageRepoInterface {
	return &MessageRepo{
		DB: db,
	}
}

func (r *MessageRepo) Save(ctx context.Context, m *message.Message) {
	r.DB.WithContext(ctx).Create(&m)
}

func (r *MessageRepo) GetByID(ctx context.Context, ID uint) (*message.Message, error) {
	var m = message.Message{ID: ID}
	err := r.DB.WithContext(ctx).First(&m).Error
	return &m, err
}

func (r *MessageRepo) LastsByRoomID(ctx context.Context, roomID uuid.UUID, c int) (messages []*message.Message, err error) {
	err = r.DB.WithContext(ctx).Where("room_id = ?", roomID).Order("time desc").Limit(c).Find(&messages).Error
	return
}

func (r *MessageRepo) LastsByTime(ctx context.Context, roomID uuid.UUID, t time.Time, c int) (messages []*message.Message, err error) {
	err = r.DB.WithContext(ctx).Where("time < ? and room_id = ?", t, roomID).Order("time desc").Limit(c).Find(&messages).Error
	return
}
