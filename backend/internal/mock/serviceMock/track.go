package servicemock

import (
	"context"

	musichub "github.com/leogues/MusicSyncHub"
)

var _ musichub.TrackService = (*TrackService)(nil)

type TrackService struct {
	FindMeTracksFn func(ctx context.Context) ([]*musichub.Track, error)
}

func (s *TrackService) FindMeTracks(ctx context.Context) ([]*musichub.Track, error) {
	return s.FindMeTracksFn(ctx)
}
