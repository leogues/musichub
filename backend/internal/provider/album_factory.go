package provider

import (
	"errors"

	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/provider/spotify/spotify_api"
)

type AlbumProviderFactoryInterface interface {
	CreateAlbumProvider(providerName string, token string) (musichub.AlbumService, error)
}

var _ AlbumProviderFactoryInterface = (*AlbumProviderFactory)(nil)

type AlbumProviderFactory struct{}

func NewAlbumProviderFactory() *AlbumProviderFactory {
	return &AlbumProviderFactory{}
}

// CreateAlbumProvider creates a new album provider
func (f *AlbumProviderFactory) CreateAlbumProvider(providerName string, token string) (musichub.AlbumService, error) {
	switch providerName {
	case musichub.AuthSourceSpotify:
		return spotify_api.NewAlbumService(token), nil
	default:
		return nil, errors.New("provider not supported")
	}
}
