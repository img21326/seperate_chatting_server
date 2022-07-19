package wait

import (
	"context"
)

type WaitRepoInterface interface {
	Add(ctx context.Context, queueName string, clientID uint)
	Len(ctx context.Context, queueName string) int
	Pop(ctx context.Context, queueName string) (clientID uint, err error)
}
