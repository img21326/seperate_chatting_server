package pair

import (
	"context"
	"fmt"
	"log"

	"github.com/img21326/fb_chat/repo/online"
	RoomRepo "github.com/img21326/fb_chat/repo/room"
	"github.com/img21326/fb_chat/repo/wait"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/structure/room"
	"github.com/img21326/fb_chat/ws/client"
)

type RedisPairUsecase struct {
	InsertClientChan <-chan *client.Client
	PairSuccessChan  chan *room.Room
	PubMessageChan   chan<- *pubmessage.PublishMessage
	WaitRepo         wait.WaitRepoInterface
	OnlineRepo       online.OnlineRepoInterface
	RoomRepo         RoomRepo.RoomRepoInterface
}

func NewRedisSubUsecase(insertClientChan <-chan *client.Client,
	pubMessageChan chan<- *pubmessage.PublishMessage,
	waitRepo wait.WaitRepoInterface, onlineRepo online.OnlineRepoInterface, roomRepo RoomRepo.RoomRepoInterface,
) PairUsecaseInterface {
	return &RedisPairUsecase{
		InsertClientChan: insertClientChan,
		PairSuccessChan:  make(chan *room.Room, 1024),
		PubMessageChan:   pubMessageChan,
		WaitRepo:         waitRepo,
		OnlineRepo:       onlineRepo,
		RoomRepo:         roomRepo,
	}
}

func (u *RedisPairUsecase) getInsertQueueName(client *client.Client) string {
	return fmt.Sprintf("%v_%v", client.User.Gender, client.WantToFind)
}

func (u *RedisPairUsecase) getPairQueueName(client *client.Client) string {
	return fmt.Sprintf("%v_%v", client.WantToFind, client.User.Gender)
}

func (u *RedisPairUsecase) Run(ctx context.Context) {
	log.Printf("[RedisPairUsecase] start")
	for {
		select {
		case client := <-u.InsertClientChan:
			// 先試著配對看看
			pairStat := false
			for {
				if u.WaitRepo.Len(ctx, u.getPairQueueName(client)) < 1 {
					break
				}
				pairClientID, err := u.WaitRepo.Pop(ctx, u.getPairQueueName(client))
				if err != nil {
					break
				}
				// 如果使用者已經下線,則在找下一個
				if !u.OnlineRepo.CheckUserOnline(ctx, pairClientID) {
					continue
				}
				room := &room.Room{
					UserId1: client.User.ID,
					UserId2: pairClientID,
					Close:   false,
				}
				u.PairSuccessChan <- room
				pairStat = true
				log.Printf("[RedisPairUsecase] pair success: %v & %v\n", room.UserId1, room.UserId2)
				break
			}
			if !pairStat {
				u.WaitRepo.Add(ctx, u.getInsertQueueName(client), client.User.ID)
				log.Printf("[RedisPairUsecase] add queue user: %v\n", client.User.ID)
			}
		case room := <-u.PairSuccessChan:
			err := u.RoomRepo.Create(room)
			var m1 *pubmessage.PublishMessage
			var m2 *pubmessage.PublishMessage
			if err != nil {
				log.Printf("create chat room err: %v", err)
				m1 = &pubmessage.PublishMessage{
					Type:     "pairError",
					SendFrom: 0,
					SendTo:   room.UserId1,
					Payload:  "create room error",
				}
				m2 = &pubmessage.PublishMessage{
					Type:     "pairError",
					SendFrom: 0,
					SendTo:   room.UserId2,
					Payload:  "create room error",
				}
			} else {
				m1 = &pubmessage.PublishMessage{
					Type:     "pairSuccess",
					SendFrom: 0,
					SendTo:   room.UserId1,
					Payload:  room.ID,
				}
				m2 = &pubmessage.PublishMessage{
					Type:     "pairSuccess",
					SendFrom: 0,
					SendTo:   room.UserId2,
					Payload:  room.ID,
				}
			}

			u.PubMessageChan <- m1
			u.PubMessageChan <- m2
			log.Printf("[RedisPairUsecase] pair user: %v %v in room: %v\n", room.UserId1, room.UserId2, room.ID)
		}
	}
}
