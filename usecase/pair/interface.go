package pair

import "context"

type PairUsecaseInterface interface {
	Run(ctx context.Context)
}
