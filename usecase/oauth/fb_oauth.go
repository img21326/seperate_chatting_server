package oauth

import (
	"context"

	"github.com/img21326/fb_chat/helper"
)

type FacebookOauthUsecase struct {
	FacebookOauth *helper.FacebookOauth
}

func NewFacebookOauthUsecase(f *helper.FacebookOauth) UsecaseOauth {
	return &FacebookOauthUsecase{
		FacebookOauth: f,
	}
}

func (f *FacebookOauthUsecase) GetLoginURL() string {
	return f.FacebookOauth.Oauth.AuthCodeURL(helper.RandString(20))
}

func (f *FacebookOauthUsecase) GetRedirectToken(ctx context.Context, code string) (*OauthToken, error) {
	token, err := f.FacebookOauth.Oauth.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return &OauthToken{
		Token:      token.AccessToken,
		ExpireTime: token.Expiry,
	}, nil
}
