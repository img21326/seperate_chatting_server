package online

import (
	"context"
	"sync"
)

type LocalOnlineRepo struct {
	ClientMap map[uint]bool
	lock      *sync.RWMutex
}

func NewLocalOnlineRepo() OnlineRepoInterface {
	return &LocalOnlineRepo{
		ClientMap: make(map[uint]bool),
		lock:      &sync.RWMutex{},
	}
}

func (r *LocalOnlineRepo) Register(ctx context.Context, clientID uint) {
	defer r.lock.Unlock()
	r.lock.Lock()
	r.ClientMap[clientID] = true
}

func (r *LocalOnlineRepo) UnRegister(ctx context.Context, clientID uint) {
	defer r.lock.Unlock()
	r.lock.Lock()
	_, ok := r.ClientMap[clientID]
	if !ok {
		return
	}
	r.ClientMap[clientID] = false
}

func (r *LocalOnlineRepo) CheckUserOnline(ctx context.Context, clientID uint) bool {
	defer r.lock.RUnlock()
	r.lock.RLock()
	status, ok := r.ClientMap[clientID]
	if !ok {
		return status
	}
	return status
}
