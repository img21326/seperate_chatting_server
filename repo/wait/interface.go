package wait

import "github.com/img21326/fb_chat/entity/ws"

type WaitRepoInterface interface {
	Add(client *ws.Client)
	Remove(client *ws.Client)
	GetFirst(*ws.Client) (*ws.Client, error)
}
