package spotify_api

import (
	"context"

	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/api"
	"github.com/leogues/MusicSyncHub/internal/provider/spotify"
	spotifyconverter "github.com/leogues/MusicSyncHub/internal/provider/spotify/spotify_converter"
)

var _ musichub.AlbumService = (*AlbumService)(nil)

type ArtistService struct {
	spotifyApi *SpotifyApi
}

func NewArtistService(token string) *ArtistService {
	return &ArtistService{
		spotifyApi: NewSpotifyApi(token),
	}
}

func (s *ArtistService) FindMeArtists(_ context.Context) ([]*musichub.Artist, error) {
	url := baseUrl + "/me/following?type=artist"

	header := s.spotifyApi.createAuthHeader()

	result, err := api.MakeAPIRequest[spotify.MeArtistsResponse](url, header)

	if err != nil {
		return nil, err
	}

	artists := spotifyconverter.ArtistsResponse(result)

	return artists, nil
}
