package online

import (
	"context"
)

type OnlineRepoInterface interface {
	Register(ctx context.Context, clientID uint)
	UnRegister(ctx context.Context, clientID uint)
	CheckUserOnline(ctx context.Context, clientID uint) bool
}
