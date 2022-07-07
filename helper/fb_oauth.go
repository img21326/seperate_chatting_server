package helper

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
	facebookOAuth "golang.org/x/oauth2/facebook"
)

type FacebookUserDetails struct {
	ID       string
	Name     string
	Email    string
	Gender   string
	Link     string
	Birthday string
}

type FacebookOauth struct {
	Oauth *oauth2.Config
	Key   string
}

func NewFacebookOauth() *FacebookOauth {
	return &FacebookOauth{
		Oauth: GetFacebookOAuthConfig(),
		Key:   RandString(20),
	}
}

func GetFacebookOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     getEnv("Facebook_ClientID", "742754863709260"),
		ClientSecret: getEnv("Facebook_ClientSecret", "4fd3867f390d4d27b5d573e902c09839"),
		RedirectURL:  getEnv("Facebook_RedirectURL", "http://localhost:8081/oauth"),
		Endpoint:     facebookOAuth.Endpoint,
		Scopes:       []string{"email", "user_gender", "user_link", "user_photos", "user_age_range"},
	}
}

func (f *FacebookOauth) GetUserInfo(token string) (FacebookUserDetails, error) {
	return GetUserInfoFromFacebook(token)
}

// GetUserInfoFromFacebook will return information of user which is fetched from facebook
func GetUserInfoFromFacebook(token string) (FacebookUserDetails, error) {
	var fbUserDetails FacebookUserDetails
	facebookUserDetailsRequest, _ := http.NewRequest("GET", "https://graph.facebook.com/me?fields=id,name,email,birthday,gender,link&access_token="+token, nil)
	facebookUserDetailsResponse, facebookUserDetailsResponseError := http.DefaultClient.Do(facebookUserDetailsRequest)

	if facebookUserDetailsResponseError != nil {
		return FacebookUserDetails{}, errors.New("Error occurred while getting information from Facebook")
	}

	decoder := json.NewDecoder(facebookUserDetailsResponse.Body)
	decoderErr := decoder.Decode(&fbUserDetails)
	defer facebookUserDetailsResponse.Body.Close()

	if decoderErr != nil {
		return FacebookUserDetails{}, errors.New("Error occurred while getting information from Facebook")
	}

	return fbUserDetails, nil
}
