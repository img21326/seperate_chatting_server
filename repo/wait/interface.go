package wait

import "github.com/img21326/fb_chat/ws/client"

type WaitRepoInterface interface {
	Add(client *client.Client)
	Remove(client *client.Client)
	GetFirst(*client.Client) (*client.Client, error)
}
