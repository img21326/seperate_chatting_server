package pair

import (
	"context"
	"fmt"
	"log"

	"github.com/img21326/fb_chat/repo/online"
	RoomRepo "github.com/img21326/fb_chat/repo/room"
	"github.com/img21326/fb_chat/repo/wait"
	errorStruct "github.com/img21326/fb_chat/structure/error"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/structure/room"
	"github.com/img21326/fb_chat/ws/client"
)

type RedisPairUsecase struct {
	InsertClientChan chan *client.Client
	PairSuccessChan  chan *room.Room
	PubMessageChan   chan<- *pubmessage.PublishMessage
	WaitRepo         wait.WaitRepoInterface
	OnlineRepo       online.OnlineRepoInterface
	RoomRepo         RoomRepo.RoomRepoInterface
}

func NewRedisSubUsecase(
	// insertClientChan <-chan *client.Client,
	// pubMessageChan chan<- *pubmessage.PublishMessage,
	waitRepo wait.WaitRepoInterface, onlineRepo online.OnlineRepoInterface, roomRepo RoomRepo.RoomRepoInterface,
) PairUsecaseInterface {
	return &RedisPairUsecase{
		InsertClientChan: make(chan *client.Client, 1024),
		PairSuccessChan:  make(chan *room.Room, 1024),
		// PubMessageChan:   pubMessageChan,
		WaitRepo:   waitRepo,
		OnlineRepo: onlineRepo,
		RoomRepo:   roomRepo,
	}
}

func (u *RedisPairUsecase) getInsertQueueName(client *client.Client) string {
	return fmt.Sprintf("%v_%v", client.User.Gender, client.WantToFind)
}

func (u *RedisPairUsecase) getPairQueueName(client *client.Client) string {
	return fmt.Sprintf("%v_%v", client.WantToFind, client.User.Gender)
}

func (u *RedisPairUsecase) SetMessageChan(messageChan chan *pubmessage.PublishMessage) {
	u.PubMessageChan = messageChan
}

func (u *RedisPairUsecase) AddToQueue(client *client.Client) {
	u.InsertClientChan <- client
}

func (u *RedisPairUsecase) TryToPair(ctx context.Context, client *client.Client) (newRoom *room.Room, err error) {
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
		newRoom = &room.Room{
			UserId1: client.User.ID,
			UserId2: pairClientID,
			Close:   false,
		}
		pairStat = true
		log.Printf("[RedisPairUsecase] pair success: %v & %v\n", newRoom.UserId1, newRoom.UserId2)
		break
	}
	if !pairStat {
		u.WaitRepo.Add(ctx, u.getInsertQueueName(client), client.User.ID)
		log.Printf("[RedisPairUsecase] add queue user: %v\n", client.User.ID)
		return nil, errorStruct.PairNotSuccess
	}
	return newRoom, nil
}
func (u *RedisPairUsecase) PairSuccess(ctx context.Context, room *room.Room) (m1 *pubmessage.PublishMessage,
	m2 *pubmessage.PublishMessage, err error) {
	err = u.RoomRepo.Create(ctx, room)
	if err != nil {
		log.Printf("create chat room err: %v", err)
		m1 = &pubmessage.PublishMessage{
			Type:     "pairError",
			SendFrom: room.UserId1,
			SendTo:   room.UserId2,
			Payload:  "create room error",
		}
		m2 = &pubmessage.PublishMessage{
			Type:     "pairError",
			SendFrom: room.UserId1,
			SendTo:   room.UserId2,
			Payload:  "create room error",
		}
	} else {
		m1 = &pubmessage.PublishMessage{
			Type:     "pairSuccess",
			SendFrom: room.UserId1,
			SendTo:   room.UserId2,
			Payload:  room.ID,
		}
		m2 = &pubmessage.PublishMessage{
			Type:     "pairSuccess",
			SendFrom: room.UserId2,
			SendTo:   room.UserId1,
			Payload:  room.ID,
		}
	}
	return
}
