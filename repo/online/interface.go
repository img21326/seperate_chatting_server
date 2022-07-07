package online

import "github.com/img21326/fb_chat/entity/ws"

type OnlineRepoInterface interface {
	Register(client *ws.Client)
	UnRegister(client *ws.Client)
	FindUserByFbID(userId uint) (*ws.Client, error)
}
