package ws

import "github.com/img21326/fb_chat/helper"

type Message struct {
	Type    string          `json:"type"`
	Message string          `json:"message"`
	Time    helper.JSONTime `json:"time"`
}
