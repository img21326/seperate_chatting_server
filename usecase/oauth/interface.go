package oauth

import (
	"context"
	"time"
)

type OauthToken struct {
	Token      string
	ExpireTime time.Time
}

type UsecaseOauth interface {
	GetLoginURL() string
	GetRedirectToken(ctx context.Context, code string) (*OauthToken, error)
}
