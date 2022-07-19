package user

import (
	"context"

	"github.com/img21326/fb_chat/structure/user"
)

type UserRepoInterFace interface {
	Create(ctx context.Context, u *user.User) error
	FindByFbID(ctx context.Context, FbId string) (*user.User, error)
}
