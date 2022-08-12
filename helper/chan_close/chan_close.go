package chan_close

import (
	"sync"

	errorStruct "github.com/img21326/fb_chat/structure/error"
)

type ChannelClose[T any] struct {
	Chan   chan T
	Closed bool
	lock   *sync.RWMutex
}

func NewChanClose[T any](c chan T) *ChannelClose[T] {
	return &ChannelClose[T]{
		Chan:   c,
		Closed: false,
		lock:   &sync.RWMutex{},
	}
}

func (c *ChannelClose[T]) Close() {
	defer c.lock.Unlock()
	c.lock.Lock()
	c.Closed = true
}

func (c *ChannelClose[T]) Push(obj T) error {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if c.Closed {
		return errorStruct.ChannelClosed
	}
	c.Chan <- obj
	return nil
}

func (c *ChannelClose[T]) Pop() T {
	return <-c.Chan
}
