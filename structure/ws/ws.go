package ws

import (
	"github.com/img21326/fb_chat/helper"
	"github.com/img21326/fb_chat/structure/room"
)

type PairSuccess struct {
	Room room.Room
}

type WebsocketMessage struct {
	Type    string          `json:"type"`
	Message string          `json:"message"`
	Time    helper.JSONTime `json:"time"`
}
