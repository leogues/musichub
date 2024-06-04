package strategy

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

func OAuth2SpotifyConfig(SpotifyClientID string, SpotifyClientSecret string) *oauth2.Config {
	redirectURL := baseUrl + "auth/spotify/callback"
	return &oauth2.Config{
		ClientID:     SpotifyClientID,
		ClientSecret: SpotifyClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"user-read-private", "user-read-email", "playlist-read-private", "playlist-read-collaborative", "user-library-read", "user-follow-read"},
		Endpoint:     spotify.Endpoint,
	}
}
