package controller

import "github.com/img21326/fb_chat/usecase/oauth"

type LoginController struct {
	OauthUsecase *oauth.UsecaseOauth
}

func NewLoginController(oauth *oauth.UsecaseOauth) *LoginController {
	return &LoginController{
		OauthUsecase: oauth,
	}
}
