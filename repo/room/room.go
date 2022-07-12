package room

import (
	"context"
	"errors"

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

func (repo *RoomRepo) Create(room *room.Room) (err error) {
	room.ID = uuid.New()
	return repo.DB.Create(room).Error
}

func (repo *RoomRepo) Close(roomId uuid.UUID) error {
	return repo.DB.Model(&room.Room{}).Where("id = ?", roomId).Update("close", true).Error
}

func (repo *RoomRepo) FindByUserId(ctx context.Context, userId uint) (room *room.Room, err error) {
	if err := repo.DB.WithContext(ctx).Order("`rooms`.`created_at` desc").Where("user_id1 = ?", userId).Or("user_id2 = ?", userId).First(&room).Error; err != nil {
		return nil, err
	}
	if room.Close {
		return nil, errors.New("RoomIsClosed")
	}
	return
}
