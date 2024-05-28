package servicemock

import (
	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/provider"
)

var _ provider.AlbumProviderFactoryInterface = (*AlbumProviderFactory)(nil)

type AlbumProviderFactory struct {
	AlbumService musichub.AlbumService
}

func (f *AlbumProviderFactory) CreateAlbumProvider(providerName string, token string) (musichub.AlbumService, error) {
	return f.AlbumService, nil
}
