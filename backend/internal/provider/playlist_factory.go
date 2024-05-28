package provider

import (
	"errors"

	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/provider/spotify/spotify_api"
)

type PlaylistProviderFactoryInterface interface {
	CreatePlaylistProvider(providerName string, token string) (musichub.PlaylistService, error)
}

var _ PlaylistProviderFactoryInterface = (*PlaylistProviderFactory)(nil)

type PlaylistProviderFactory struct{}

func NewPlaylistProviderFactory() *PlaylistProviderFactory {
	return &PlaylistProviderFactory{}
}

// CreatePlaylistProvider creates a new playlist provider
func (f *PlaylistProviderFactory) CreatePlaylistProvider(providerName string, token string) (musichub.PlaylistService, error) {
	switch providerName {
	case musichub.AuthSourceSpotify:
		return spotify_api.NewPlaylistService(token), nil
	default:
		return nil, errors.New("provider not supported")
	}
}
