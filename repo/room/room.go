package room

import (
	"errors"
	"log"

	"github.com/google/uuid"
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

func (repo *RoomRepo) Create(room *Room) (err error) {
	room.ID = uuid.New()
	return repo.DB.Create(room).Error
}

func (repo *RoomRepo) Close(roomId uuid.UUID) error {
	return repo.DB.Model(&Room{}).Where("id = ?", roomId).Update("close", true).Error
}

func (repo *RoomRepo) FindByUserId(userId uint) (room *Room, err error) {
	if err := repo.DB.Where("user_id1 = ?", userId).Or("user_id2 = ?", userId).Last(&room).Error; err != nil {
		return nil, err
	}
	log.Printf("room: %+v", room)
	if room.Close {
		return nil, errors.New("RoomIsClosed")
	}
	return
}
