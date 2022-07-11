package ws

import "context"

type WebsocketUsecaseInterface interface {
	Run(ctx context.Context)
}
