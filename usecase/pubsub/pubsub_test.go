package pubsub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/img21326/fb_chat/mock"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/stretchr/testify/assert"
)

func TestSubscribe(t *testing.T) {
	c := gomock.NewController(t)
	pubsubRepo := mock.NewMockPubSubRepoInterface(c)
	pubsubRepo.EXPECT().Sub(gomock.Any(), "test").Times(1).Return(func() ([]byte, error) {
		return []byte("test"), nil
	})

	pubsubUsecase := RedisSubUsecase{PubSubRepo: pubsubRepo}
	ctx := context.Background()
	sub := pubsubUsecase.Subscribe(ctx, "test")

	msg, err := sub()

	assert.Equal(t, msg, []byte("test"))
	assert.Nil(t, err)
}

func TestPublish(t *testing.T) {
	c := gomock.NewController(t)
	pubsubRepo := mock.NewMockPubSubRepoInterface(c)

	m := &pubmessage.PublishMessage{
		Payload: "test",
	}
	jsonMessage, _ := json.Marshal(m)
	pubsubRepo.EXPECT().Pub(gomock.Any(), "test", jsonMessage).Times(1)

	pubsubUsecase := RedisSubUsecase{PubSubRepo: pubsubRepo}
	ctx := context.Background()
	pubsubUsecase.Publish(ctx, "test", m)
}
