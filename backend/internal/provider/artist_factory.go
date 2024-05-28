package provider

import (
	"errors"

	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/provider/spotify/spotify_api"
)

type ArtistProviderFactoryInterface interface {
	CreateArtistProvider(providerName string, token string) (musichub.ArtistService, error)
}

var _ ArtistProviderFactoryInterface = (*ArtistProviderFactory)(nil)

type ArtistProviderFactory struct{}

func NewArtistProviderFactory() *ArtistProviderFactory {
	return &ArtistProviderFactory{}
}

// CreateArtistProvider creates a new artist provider
func (f *ArtistProviderFactory) CreateArtistProvider(providerName string, token string) (musichub.ArtistService, error) {
	switch providerName {
	case musichub.AuthSourceSpotify:
		return spotify_api.NewArtistService(token), nil
	default:
		return nil, errors.New("provider not supported")
	}
}
