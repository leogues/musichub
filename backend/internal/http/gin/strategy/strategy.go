package strategy

import (
	"os"

	"golang.org/x/oauth2"
)

var (
	baseUrl string
)

func init() {
	baseUrl = os.Getenv("BASE_URL")
}

type OAuthStrategys struct {
	Google  *oauth2.Config
	Spotify *oauth2.Config
}

func NewOAuthStrategys() *OAuthStrategys {
	return &OAuthStrategys{}
}
