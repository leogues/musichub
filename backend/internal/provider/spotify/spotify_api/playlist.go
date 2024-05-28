package spotify_api

import (
	"context"

	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/api"
	"github.com/leogues/MusicSyncHub/internal/provider/spotify"
	spotifyconverter "github.com/leogues/MusicSyncHub/internal/provider/spotify/spotify_converter"
)

var _ musichub.PlaylistService = (*PlaylistService)(nil)

type PlaylistService struct {
	SpotifyApi *SpotifyApi
}

func NewPlaylistService(token string) *PlaylistService {
	return &PlaylistService{
		SpotifyApi: NewSpotifyApi(token),
	}
}

func (s *PlaylistService) FindMePlaylists(_ context.Context) ([]*musichub.Playlist, error) {
	url := baseUrl + "/me/playlists"

	header := s.SpotifyApi.createAuthHeader()

	result, err := api.MakeAPIRequest[spotify.MePlaylistsResponse](url, header)

	if err != nil {
		return nil, err
	}

	playlists := spotifyconverter.MePlaylistsResponse(result)

	return playlists, nil
}

func (s *PlaylistService) FindPlaylistByID(_ context.Context, id string) (*musichub.Playlist, error) {
	url := baseUrl + "/playlists/" + id

	header := s.SpotifyApi.createAuthHeader()

	result, err := api.MakeAPIRequest[spotify.PlaylistResponse](url, header)

	if err != nil {
		return nil, err
	}

	playlist := spotifyconverter.PlaylistResponse(result)

	return playlist, nil
}
