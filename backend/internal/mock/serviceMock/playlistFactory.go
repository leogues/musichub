package servicemock

import (
	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/provider"
)

var _ provider.PlaylistProviderFactoryInterface = (*PlaylistProviderFactory)(nil)

type PlaylistProviderFactory struct {
	PlaylistService musichub.PlaylistService
}

func (f *PlaylistProviderFactory) CreatePlaylistProvider(providerName string, token string) (musichub.PlaylistService, error) {
	return f.PlaylistService, nil
}
