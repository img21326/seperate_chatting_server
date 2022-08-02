package hub

import (
	"context"
	"log"
	"time"

	MessageHub "github.com/img21326/fb_chat/hub/message"
	PairHub "github.com/img21326/fb_chat/hub/pair"
	PubSubHub "github.com/img21326/fb_chat/hub/pubsub"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/usecase/message"
	"github.com/img21326/fb_chat/usecase/pair"
	"github.com/img21326/fb_chat/usecase/pubsub"
	"github.com/img21326/fb_chat/usecase/ws"
	"github.com/img21326/fb_chat/ws/client"
)

func StartHub(
	ctx context.Context,
	pubsubUsecase pubsub.SubMessageUsecaseInterface,
	pairUsecase pair.PairUsecaseInterface,
	messageUsecase message.MessageUsecaseInterface,
	wsUsecase ws.WebsocketUsecaseInterface,
) (
	pubMessageChan chan *pubmessage.PublishMessage,
	clientQueueChan chan *client.Client,
) {
	pubMessageChan = make(chan *pubmessage.PublishMessage, 4096)
	clientQueueChan = make(chan *client.Client, 4096)
	subMessageChan := make(chan *pubmessage.PublishMessage, 4096)

	subHub := PubSubHub.NewSubHub(pubsubUsecase)
	pubHub := PubSubHub.NewPubHub(pubsubUsecase)
	messageHub := MessageHub.NewMessageHub(messageUsecase, wsUsecase, subMessageChan)
	pairHub := PairHub.NewPairHub(pairUsecase, pubMessageChan, clientQueueChan)

	subCtx, cancel := context.WithCancel(context.Background())

	go subHub.Run(subCtx, "message", subMessageChan)
	go pubHub.Run(ctx, "message", pubMessageChan)
	go messageHub.Run(ctx, 5*time.Second)
	go pairHub.Run(ctx)

	// 最後送出的訊息 要收到並處理完
	go func(ctx context.Context, cancel context.CancelFunc) {
		<-ctx.Done()
		time.Sleep(time.Duration(5 * time.Second))
		cancel()
		log.Printf("[Hub] close sub")
	}(ctx, cancel)

	return pubMessageChan, clientQueueChan
}
