package room

import (
	"context"

	"github.com/google/uuid"
	"github.com/img21326/fb_chat/structure/room"
	"gorm.io/gorm"
)

type RoomRepo struct {
	DB *gorm.DB
}

func NewRoomRepo(db *gorm.DB) RoomRepoInterface {
	return &RoomRepo{
		DB: db,
	}
}

func (repo *RoomRepo) Create(ctx context.Context, room *room.Room) (err error) {
	room.UUID = uuid.New()
	return repo.DB.WithContext(ctx).Create(room).Error
}

func (repo *RoomRepo) Close(ctx context.Context, roomId uuid.UUID) error {
	return repo.DB.WithContext(ctx).Model(&room.Room{}).Where("UUID = ?", roomId).Update("close", true).Error
}

func (repo *RoomRepo) FindByUserId(ctx context.Context, userId uint) (room *room.Room, err error) {
	if err := repo.DB.WithContext(ctx).Order("`rooms`.`created_at` desc").Where("user_id1 = ?", userId).Or("user_id2 = ?", userId).First(&room).Error; err != nil {
		return nil, err
	}
	return
}
