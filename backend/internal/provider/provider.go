package provider

import musichub "github.com/leogues/MusicSyncHub"

var _ musichub.ProviderService = (*ProviderService)(nil)

type ProviderService struct{}

func NewProviderService() *ProviderService {
	return &ProviderService{}
}

// Platforms returns the list of platforms available
func (s *ProviderService) Platforms() []*musichub.Platform {
	return []*musichub.Platform{
		{
			Id:           "0",
			TypePlatform: "music",
			Name:         "spotify",
			Label:        "Spotify",
		},
		{
			Id:           "1",
			TypePlatform: "music",
			Name:         "google",
			Label:        "Google",
		},
	}
}
