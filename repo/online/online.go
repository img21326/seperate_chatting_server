package online

import (
	"errors"
	"sync"

	"github.com/img21326/fb_chat/entity/ws"
	"gorm.io/gorm"
)

type OnlineRepo struct {
	DB        *gorm.DB
	ClientMap map[uint]*ws.Client
	lock      *sync.Mutex
}

func NewOnlineRepo(db *gorm.DB) OnlineRepoInterface {
	return &OnlineRepo{
		ClientMap: make(map[uint]*ws.Client),
		lock:      &sync.Mutex{},
	}
}

func (r *OnlineRepo) Register(client *ws.Client) {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()
	r.ClientMap[client.User.ID] = client
}

func (r *OnlineRepo) UnRegister(client *ws.Client) {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()
	if _, ok := r.ClientMap[client.User.ID]; ok {
		delete(r.ClientMap, client.User.ID)
	}
}

func (r *OnlineRepo) FindUserByFbID(userId uint) (*ws.Client, error) {
	if client, ok := r.ClientMap[userId]; ok {
		return client, nil
	} else {
		return nil, errors.New("RecordNotFound")
	}
}
