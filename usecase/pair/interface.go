package pair

import (
	"context"

	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/structure/room"
	"github.com/img21326/fb_chat/ws/client"
)

type PairUsecaseInterface interface {
	TryToPair(ctx context.Context, client *client.Client) (room *room.Room, err error)
	PairSuccess(ctx context.Context, room *room.Room) (m1 *pubmessage.PublishMessage, m2 *pubmessage.PublishMessage, err error)
}
