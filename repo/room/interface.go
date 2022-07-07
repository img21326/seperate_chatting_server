package room

import "github.com/google/uuid"

type Room struct {
	ID      uuid.UUID
	UserId1 uint
	UserId2 uint
	Close   bool
}

type RoomRepoInterface interface {
	Create(room *Room) error
	Close(roomId uuid.UUID)
	FindByUserId(userId uint) (*Room, error)
}
