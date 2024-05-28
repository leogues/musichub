package servicemock

import (
	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/provider"
)

var _ provider.TrackProviderFactoryInterface = (*TrackProviderFactory)(nil)

type TrackProviderFactory struct {
	TrackService musichub.TrackService
}

func (f *TrackProviderFactory) CreateTrackProvider(providerName string, token string) (musichub.TrackService, error) {
	return f.TrackService, nil
}
