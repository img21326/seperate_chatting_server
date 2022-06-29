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

type UsecaseOauthInterFace interface {
	GetLoginURL() string
	GetRedirectToken(ctx context.Context, key string, code string) (*OauthToken, error)
	GetUser(token string) (*OauthUser, error)
}
