package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	facebookOAuth "golang.org/x/oauth2/facebook"
)

type FacebookUserDetails struct {
	ID     string
	Name   string
	Email  string
	Gender string
	Link   string
}

func GetFacebookOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     "742754863709260",
		ClientSecret: "4fd3867f390d4d27b5d573e902c09839",
		RedirectURL:  "https://b6cc-61-223-238-240.ngrok.io/oauth",
		Endpoint:     facebookOAuth.Endpoint,
		Scopes:       []string{"email", "user_gender", "user_link", "user_photos", "user_birthday"},
	}
}

// GetRandomOAuthStateString will return random string
func GetRandomOAuthStateString() string {
	return "SomeRandomStringAlgorithmForMoreSecurity"
}

// GetUserInfoFromFacebook will return information of user which is fetched from facebook
func GetUserInfoFromFacebook(token string) (FacebookUserDetails, error) {
	var fbUserDetails FacebookUserDetails
	facebookUserDetailsRequest, _ := http.NewRequest("GET", "https://graph.facebook.com/me?fields=id,name,email&access_token="+token, nil)
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

type UserDetails struct {
	Name     string
	Email    string
	Password string
}

// FacebookUserDetails is struct used for user details

func main() {
	server := gin.Default()
	server.GET("/login", func(ctx *gin.Context) {
		var OAuth2Config = GetFacebookOAuthConfig()
		url := OAuth2Config.AuthCodeURL(GetRandomOAuthStateString())
		ctx.Redirect(http.StatusFound, url)
	})
	server.GET("/oauth", func(ctx *gin.Context) {
		code := ctx.Request.FormValue("code")
		var OAuth2Config = GetFacebookOAuthConfig()

		token, _ := OAuth2Config.Exchange(oauth2.NoContext, code)
		ctx.JSON(200, gin.H{
			"token": token,
		})
	})
	server.Run(":8081")
}
