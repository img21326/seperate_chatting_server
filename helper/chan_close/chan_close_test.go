package chan_close

import (
	"testing"

	errorStruct "github.com/img21326/fb_chat/structure/error"
	"github.com/stretchr/testify/assert"
)

func TestChanClosePush(t *testing.T) {
	chan_ := make(chan interface{}, 8)
	c := NewChanClose(chan_)
	err := c.Push("abc")
	assert.Nil(t, err)
}

func TestChanClosePushWithClose(t *testing.T) {
	chan_ := make(chan interface{}, 8)
	c := NewChanClose(chan_)
	c.Close()
	err := c.Push("abc")
	assert.Equal(t, err, errorStruct.ChannelClosed)
}

func TestChanClosePop(t *testing.T) {
	chan_ := make(chan interface{}, 8)
	c := NewChanClose(chan_)
	err := c.Push("abc")
	assert.Nil(t, err)

	get := c.Pop()
	assert.Equal(t, get.(string), "abc")
}
