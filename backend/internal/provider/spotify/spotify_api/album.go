package spotify_api

import (
	"context"

	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/api"
	"github.com/leogues/MusicSyncHub/internal/provider/spotify"
	spotifyconverter "github.com/leogues/MusicSyncHub/internal/provider/spotify/spotify_converter"
)

var _ musichub.AlbumService = (*AlbumService)(nil)

type AlbumService struct {
	spotifyApi *SpotifyApi
}

func NewAlbumService(token string) *AlbumService {
	return &AlbumService{
		spotifyApi: NewSpotifyApi(token),
	}
}

func (s *AlbumService) FindMeAlbums(_ context.Context) ([]*musichub.Album, error) {
	url := baseUrl + "/me/albums"

	header := s.spotifyApi.createAuthHeader()

	result, err := api.MakeAPIRequest[spotify.MeAlbumsResponse](url, header)

	if err != nil {
		return nil, err
	}

	albums := spotifyconverter.MeAlbumsResponse(result)

	return albums, nil
}

func (s *AlbumService) FindAlbumByID(_ context.Context, id string) (*musichub.Album, error) {
	url := baseUrl + "/albums/" + id

	header := s.spotifyApi.createAuthHeader()

	result, err := api.MakeAPIRequest[spotify.AlbumResponse](url, header)

	if err != nil {
		return nil, err
	}

	album := spotifyconverter.AlbumResponse(result)

	return album, nil
}
