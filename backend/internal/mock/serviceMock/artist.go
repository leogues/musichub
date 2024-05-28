package servicemock

import (
	"context"

	musichub "github.com/leogues/MusicSyncHub"
)

var _ musichub.ArtistService = (*ArtistService)(nil)

type ArtistService struct {
	FindMeArtistsFn func(ctx context.Context) ([]*musichub.Artist, error)
}

func (s *ArtistService) FindMeArtists(ctx context.Context) ([]*musichub.Artist, error) {
	return s.FindMeArtistsFn(ctx)
}
