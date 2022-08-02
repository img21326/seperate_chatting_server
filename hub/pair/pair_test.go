package pair

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/img21326/fb_chat/mock"
	errorStruct "github.com/img21326/fb_chat/structure/error"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/structure/room"
	"github.com/img21326/fb_chat/ws/client"
	"github.com/stretchr/testify/assert"
)

func TestInsertClientPairSuccess(t *testing.T) {
	c := gomock.NewController(t)

	cli := &client.Client{}
	cli.User.ID = 1

	r := &room.Room{
		UUID: uuid.New(),
	}

	ctx, cancel := context.WithCancel(context.Background())

	pairUsecase := mock.NewMockPairUsecaseInterface(c)
	pairUsecase.EXPECT().TryToPair(gomock.Any(), cli).Times(1).DoAndReturn(
		func(ctx context.Context, client *client.Client) (*room.Room, error) {
			cancel()
			return r, nil
		})

	pairHub := &PairHub{
		PairUsecase:      pairUsecase,
		InsertClientChan: make(chan *client.Client, 1),
		PairSuccessChan:  make(chan *room.Room, 1),
	}

	go pairHub.Run(ctx)
	pairHub.InsertClientChan <- cli

	assert.Equal(t, r, <-pairHub.PairSuccessChan)
}

func TestInsertClientPairNotSuccess(t *testing.T) {
	c := gomock.NewController(t)

	cli := &client.Client{}
	cli.User.ID = 1

	ctx, cancel := context.WithCancel(context.Background())

	pairUsecase := mock.NewMockPairUsecaseInterface(c)
	pairUsecase.EXPECT().TryToPair(gomock.Any(), cli).Times(1).DoAndReturn(
		func(ctx context.Context, client *client.Client) (*room.Room, error) {
			cancel()
			return nil, errorStruct.PairNotSuccess
		})
	pairUsecase.EXPECT().AddToQueue(gomock.Any(), cli).Times(1)

	pairHub := &PairHub{
		PairUsecase:      pairUsecase,
		InsertClientChan: make(chan *client.Client, 1),
		PairSuccessChan:  make(chan *room.Room, 1),
	}

	pairHub.InsertClientChan <- cli
	pairHub.Run(ctx)
}

func TestPairSuccess(t *testing.T) {
	c := gomock.NewController(t)

	r := &room.Room{UUID: uuid.New()}

	pub1 := &pubmessage.PublishMessage{
		SendFrom: 1,
		SendTo:   2,
	}
	pub2 := &pubmessage.PublishMessage{
		SendFrom: 2,
		SendTo:   1,
	}

	ctx, cancel := context.WithCancel(context.Background())

	pairUsecase := mock.NewMockPairUsecaseInterface(c)
	pairUsecase.EXPECT().PairSuccess(gomock.Any(), r).Times(1).DoAndReturn(
		func(ctx context.Context, room *room.Room) (*pubmessage.PublishMessage, *pubmessage.PublishMessage, error) {
			cancel()
			return pub1, pub2, nil
		})

	pairHub := &PairHub{
		PairUsecase:      pairUsecase,
		PubMessageChan:   make(chan *pubmessage.PublishMessage, 2),
		PairSuccessChan:  make(chan *room.Room, 1),
		InsertClientChan: make(chan *client.Client),
	}

	pairHub.PairSuccessChan <- r
	pairHub.Run(ctx)

	assert.Equal(t, pub1, <-pairHub.PubMessageChan)
	assert.Equal(t, pub2, <-pairHub.PubMessageChan)
}

func TestGracefulShutdown(t *testing.T) {
	c := gomock.NewController(t)

	ctx, cancel := context.WithCancel(context.Background())
	pairUsecase := mock.NewMockPairUsecaseInterface(c)

	pub1 := &pubmessage.PublishMessage{
		SendFrom: 1,
		SendTo:   2,
	}
	pub2 := &pubmessage.PublishMessage{
		SendFrom: 2,
		SendTo:   1,
	}
	pairUsecase.EXPECT().PairSuccess(gomock.Any(), gomock.Any()).Times(30).DoAndReturn(
		func(ctx context.Context, room *room.Room) (*pubmessage.PublishMessage, *pubmessage.PublishMessage, error) {
			time.Sleep(time.Microsecond * 600)
			return pub1, pub2, nil
		})

	pairHub := &PairHub{
		PairUsecase:      pairUsecase,
		PubMessageChan:   make(chan *pubmessage.PublishMessage, 60),
		PairSuccessChan:  make(chan *room.Room, 30),
		InsertClientChan: make(chan *client.Client),
	}
	go pairHub.Run(ctx)

	for i := 1; i <= 30; i++ {
		pairHub.PairSuccessChan <- &room.Room{}
	}
	cancel()
	time.Sleep(3 * time.Second)
}
