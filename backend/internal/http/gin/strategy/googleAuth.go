package strategy

import (
	"net/http"

	"github.com/leogues/MusicSyncHub/api"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	googleBaseUrl = "https://www.googleapis.com/oauth2/v1"
)

type clientGoogle struct {
	token *oauth2.Token
}

type UserInfo struct {
	ID         *string `json:"id"`
	Name       *string `json:"given_name"`
	FamilyName *string `json:"family_name"`
	Email      *string `json:"email"`
	Picture    *string `json:"picture"`
}

func NewClientGoogle(token *oauth2.Token) *clientGoogle {
	return &clientGoogle{
		token: token,
	}
}

func OAuth2GoogleConfig(GoogleClientID string, GoogleClientSecret string) *oauth2.Config {
	redirectURL := baseUrl + "auth/google/callback"
	return &oauth2.Config{
		ClientID:     GoogleClientID,
		ClientSecret: GoogleClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func (client *clientGoogle) GoogleUserInfo() (*UserInfo, error) {
	url := googleBaseUrl + "/userinfo"

	header := http.Header{}
	header.Add("Authorization", "Bearer "+client.token.AccessToken)

	userInfo, err := api.MakeAPIRequest[*UserInfo](url, header)

	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
