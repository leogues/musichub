package strategy

import "golang.org/x/oauth2"

type OAuthStrategys struct {
	Google  *oauth2.Config
	Spotify *oauth2.Config
}

func NewOAuthStrategys() *OAuthStrategys {
	return &OAuthStrategys{}
}
