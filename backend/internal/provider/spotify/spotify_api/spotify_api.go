package spotify_api

import (
	"net/http"
)

const (
	baseUrl = "https://api.spotify.com/v1"
)

type SpotifyApi struct {
	token string
}

func NewSpotifyApi(token string) *SpotifyApi {
	return &SpotifyApi{
		token: token,
	}
}

func (a *SpotifyApi) createAuthHeader() http.Header {
	header := http.Header{}
	header.Set("Authorization", "Bearer "+a.token)
	return header
}
