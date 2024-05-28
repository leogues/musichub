package servicemock

import musichub "github.com/leogues/MusicSyncHub"

var _ musichub.ProviderService = (*ProviderService)(nil)

type ProviderService struct {
	PlatformsFn func() []*musichub.Platform
}

func (s *ProviderService) Platforms() []*musichub.Platform {
	return s.PlatformsFn()
}
