package spotify_api

import (
	"context"

	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/api"
	"github.com/leogues/MusicSyncHub/internal/provider/spotify"
	spotifyconverter "github.com/leogues/MusicSyncHub/internal/provider/spotify/spotify_converter"
)

var _ musichub.TrackService = (*TrackService)(nil)

type TrackService struct {
	spotifyApi *SpotifyApi
}

func NewTrackService(token string) *TrackService {
	return &TrackService{
		spotifyApi: NewSpotifyApi(token),
	}
}

func (s *TrackService) FindMeTracks(_ context.Context) ([]*musichub.Track, error) {
	url := baseUrl + "/me/tracks"

	header := s.spotifyApi.createAuthHeader()

	result, err := api.MakeAPIRequest[spotify.MeTracksResponse](url, header)

	if err != nil {
		return nil, err
	}

	tracks := spotifyconverter.MeTracksResponse(result)

	return tracks, nil
}
