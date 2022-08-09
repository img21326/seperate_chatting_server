package chan_close

import (
	"sync"

	errorStruct "github.com/img21326/fb_chat/structure/error"
)

type ChannelClose struct {
	Chan   chan interface{}
	Closed bool
	lock   *sync.RWMutex
}

func NewChanClose(c chan interface{}) *ChannelClose {
	return &ChannelClose{
		Chan:   c,
		Closed: false,
		lock:   &sync.RWMutex{},
	}
}

func (c *ChannelClose) Close() {
	defer c.lock.Unlock()
	c.lock.Lock()
	c.Closed = true
}

func (c *ChannelClose) Push(obj interface{}) error {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if c.Closed {
		return errorStruct.ChannelClosed
	}
	c.Chan <- obj
	return nil
}

func (c *ChannelClose) Pop() interface{} {
	return <-c.Chan
}
