package online

import (
	"errors"
	"sync"

	"github.com/img21326/fb_chat/ws/client"
)

type OnlineRepo struct {
	ClientMap map[uint]*client.Client
	lock      *sync.Mutex
}

func NewOnlineRepo() OnlineRepoInterface {
	return &OnlineRepo{
		ClientMap: make(map[uint]*client.Client),
		lock:      &sync.Mutex{},
	}
}

func (r *OnlineRepo) Register(client *client.Client) {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()
	r.ClientMap[client.User.ID] = client
}

func (r *OnlineRepo) UnRegister(client *client.Client) {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()
	if _, ok := r.ClientMap[client.User.ID]; ok {
		delete(r.ClientMap, client.User.ID)
	}
}

func (r *OnlineRepo) FindUserByFbID(userId uint) (*client.Client, error) {
	if client, ok := r.ClientMap[userId]; ok {
		return client, nil
	} else {
		return nil, errors.New("RecordNotFound")
	}
}
