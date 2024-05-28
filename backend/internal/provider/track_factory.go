package provider

import (
	"errors"

	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/provider/spotify/spotify_api"
)

type TrackProviderFactoryInterface interface {
	CreateTrackProvider(providerName string, token string) (musichub.TrackService, error)
}

var _ TrackProviderFactoryInterface = (*TrackProviderFactory)(nil)

type TrackProviderFactory struct{}

func NewTrackProviderFactory() *TrackProviderFactory {
	return &TrackProviderFactory{}
}

// CreateTrackProvider creates a new track provider
func (f *TrackProviderFactory) CreateTrackProvider(providerName string, token string) (musichub.TrackService, error) {
	switch providerName {
	case musichub.AuthSourceSpotify:
		return spotify_api.NewTrackService(token), nil
	default:
		return nil, errors.New("provider not supported")
	}
}
