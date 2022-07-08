package online

import "github.com/img21326/fb_chat/ws/client"

type OnlineRepoInterface interface {
	Register(client *client.Client)
	UnRegister(client *client.Client)
	FindUserByFbID(userId uint) (*client.Client, error)
}
