package room

import (
	"github.com/google/uuid"
	"github.com/img21326/fb_chat/structure/room"
)

type RoomRepoInterface interface {
	Create(room *room.Room) error
	Close(roomId uuid.UUID) error
	FindByUserId(userId uint) (*room.Room, error)
}
