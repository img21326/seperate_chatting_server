package room

import (
	"context"

	"github.com/google/uuid"
	"github.com/img21326/fb_chat/structure/room"
)

type RoomRepoInterface interface {
	Create(ctx context.Context, room *room.Room) error
	Close(ctx context.Context, roomId uuid.UUID) error
	FindByUserId(ctx context.Context, userId uint) (*room.Room, error)
}
