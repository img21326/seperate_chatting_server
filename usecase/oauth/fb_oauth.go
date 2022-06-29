package oauth

import (
	"context"
	"time"

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

func (f *FacebookOauthUsecase) GetUser(token string) (*OauthUser, error) {
	user, err := f.FacebookOauth.GetUserInfo(token)
	if err != nil {
		return nil, err
	}
	birth, err := time.Parse("01/02/2006", user.Birth)
	if err != nil {
		return nil, err
	}
	return &OauthUser{
		ID:     user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Gender: user.Gender,
		Link:   user.Link,
		Birth:  birth,
	}, nil
}
