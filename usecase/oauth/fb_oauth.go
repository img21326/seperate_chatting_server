package oauth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/img21326/fb_chat/helper"
)

type FacebookOauthUsecase struct {
	FacebookOauth *helper.FacebookOauth
}

func NewFacebookOauthUsecase(f *helper.FacebookOauth) UsecaseOauthInterFace {
	return &FacebookOauthUsecase{
		FacebookOauth: f,
	}
}

func (f *FacebookOauthUsecase) GetLoginURL() string {
	return f.FacebookOauth.Oauth.AuthCodeURL(f.FacebookOauth.Key)
}

func (f *FacebookOauthUsecase) GetRedirectToken(ctx context.Context, key string, code string) (*OauthToken, error) {
	if key != f.FacebookOauth.Key {
		return nil, errors.New("Except Key")
	}
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
	fmt.Printf("%+v", user)
	if err != nil {
		return nil, err
	}
	birth, err := time.Parse("01/02/2006", user.Birthday)
	if err != nil {
		fmt.Printf("%+v", err)
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
