package pair

import (
	"context"
	"log"

	errorStruct "github.com/img21326/fb_chat/structure/error"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/structure/room"
	PairUsecase "github.com/img21326/fb_chat/usecase/pair"
	"github.com/img21326/fb_chat/ws/client"
)

type PairHub struct {
	PubMessageChan   chan *pubmessage.PublishMessage
	InsertClientChan chan *client.Client
	PairSuccessChan  chan *room.Room
	PairUsecase      PairUsecase.PairUsecaseInterface
}

func NewPairHub(pairUsecase PairUsecase.PairUsecaseInterface,
	pubMessageChan chan *pubmessage.PublishMessage,
	insertClientChan chan *client.Client) *PairHub {
	return &PairHub{
		InsertClientChan: insertClientChan,
		PubMessageChan:   pubMessageChan,
		PairSuccessChan:  make(chan *room.Room, 1024),
		PairUsecase:      pairUsecase,
	}
}

func (h *PairHub) Run(ctx context.Context) {
	log.Printf("[PairHub] start")
	if h.PubMessageChan == nil {
		log.Printf("[PairHub] not set message chan")
	}
	for {
		select {
		case <-ctx.Done():
			return
		case client := <-h.InsertClientChan:
			// 先試著配對看看
			c := context.Background()
			room, err := h.PairUsecase.TryToPair(c, client)
			if err != nil {
				if err == errorStruct.PairNotSuccess {
					h.PairUsecase.AddToQueue(c, client)
				}
			}
			h.PairSuccessChan <- room
		case room := <-h.PairSuccessChan:
			c := context.Background()
			m1, m2, err := h.PairUsecase.PairSuccess(c, room)
			if err != nil {
				log.Printf("[RedisPairHub] error: %+v, room: %+v", err, room)
			}
			h.PubMessageChan <- m1
			h.PubMessageChan <- m2
			log.Printf("[RedisPairUsecase] pair user: %v %v in room: %v\n", room.UserId1, room.UserId2, room.ID)
		}
	}
}
