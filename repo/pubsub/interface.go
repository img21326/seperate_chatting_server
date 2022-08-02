package pubsub

import (
	"context"
)

type PubSubRepoInterface interface {
	Sub(ctx context.Context, topic string) func() ([]byte, error)
	Pub(ctx context.Context, topic string, message []byte) error
}
