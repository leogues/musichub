package servicemock

import (
	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/provider"
)

var _ provider.ArtistProviderFactoryInterface = (*ArtistProviderFactory)(nil)

type ArtistProviderFactory struct {
	ArtistService musichub.ArtistService
}

func (f *ArtistProviderFactory) CreateArtistProvider(providerName string, token string) (musichub.ArtistService, error) {
	return f.ArtistService, nil
}
