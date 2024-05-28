package servicemock

import (
	"context"

	musichub "github.com/leogues/MusicSyncHub"
)

var _ musichub.PlaylistService = (*PlaylistService)(nil)

type PlaylistService struct {
	FindMePlaylistsFn  func(ctx context.Context) ([]*musichub.Playlist, error)
	FindPlaylistByIDFn func(ctx context.Context, id string) (*musichub.Playlist, error)
}

func (s *PlaylistService) FindMePlaylists(ctx context.Context) ([]*musichub.Playlist, error) {
	return s.FindMePlaylistsFn(ctx)
}

func (s *PlaylistService) FindPlaylistByID(ctx context.Context, id string) (*musichub.Playlist, error) {
	return s.FindPlaylistByIDFn(ctx, id)
}
