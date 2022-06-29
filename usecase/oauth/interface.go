package oauth

import (
	"context"
	"time"
)

type OauthToken struct {
	Token      string
	ExpireTime time.Time
}

type OauthUser struct {
	ID     string
	Name   string
	Email  string
	Gender string
	Link   string
	Birth  time.Time
}

type UsecaseOauth interface {
	GetLoginURL() string
	GetRedirectToken(ctx context.Context, code string) (*OauthToken, error)
	GetUser(token string) (*OauthUser, error)
}
