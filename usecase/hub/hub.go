package hub

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/repo/online"
	"github.com/img21326/fb_chat/repo/room"
	"github.com/img21326/fb_chat/repo/user"
	"github.com/img21326/fb_chat/repo/wait"
	"github.com/img21326/fb_chat/ws/client"
)

type HubUsecase struct {
	UserRepo    user.UserRepoInterFace
	MessageRepo message.MessageRepoInterface
	OnlineRepo  online.OnlineRepoInterface
	RoomRepo    room.RoomRepoInterface
	WaitRepo    wait.WaitRepoInterface
}

func NewHubUsecase(userRepo user.UserRepoInterFace, messageRepo message.MessageRepoInterface,
	onlineRepo online.OnlineRepoInterface, roomRepo room.RoomRepoInterface, waitRepo wait.WaitRepoInterface,
) HubUsecaseInterface {
	return &HubUsecase{
		UserRepo:    userRepo,
		MessageRepo: messageRepo,
		OnlineRepo:  onlineRepo,
		RoomRepo:    roomRepo,
		WaitRepo:    waitRepo,
	}
}

func (u *HubUsecase) RegisterOnline(client *client.Client) {
	u.OnlineRepo.Register(client)
}

func (u *HubUsecase) UnRegisterOnline(client *client.Client) {
	u.OnlineRepo.UnRegister(client)
}

func (u *HubUsecase) FindOnlineUserByUserID(userId uint) (*client.Client, error) {
	return u.OnlineRepo.FindUserByFbID(userId)
}

func (u *HubUsecase) FindUserByFbID(FbId string) (*user.UserModel, error) {
	return u.UserRepo.FindByFbID(FbId)
}

func (u *HubUsecase) GetFirstQueueUser(client *client.Client) (*client.Client, error) {
	return u.WaitRepo.GetFirst(client)
}

func (u *HubUsecase) AddUserToQueue(client *client.Client) {
	u.WaitRepo.Add(client)
}

func (u *HubUsecase) DeleteuserFromQueue(client *client.Client) {
	u.WaitRepo.Remove(client)
}

func (u *HubUsecase) CreateRoom(room *room.Room) error {
	return u.RoomRepo.Create(room)
}

func (u *HubUsecase) CloseRoom(uuid uuid.UUID) error {
	return u.RoomRepo.Close(uuid)
}

func (u *HubUsecase) FindRoomByUserId(userId uint) (*room.Room, error) {
	return u.RoomRepo.FindByUserId(userId)
}

func (u *HubUsecase) SaveMesssage(ctx context.Context, message *message.MessageModel) {
	u.MessageRepo.Save(ctx, message)
}

func (u *HubUsecase) SendMessage(ctx context.Context, message message.PublishMessage) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return u.MessageRepo.Publish(ctx, jsonMessage)
}
