package servicemock

import (
	"context"

	musichub "github.com/leogues/MusicSyncHub"
)

var _ musichub.AlbumService = (*AlbumService)(nil)

type AlbumService struct {
	FindMeAlbumsFn    func(ctx context.Context) ([]*musichub.Album, error)
	FindMeAlbumByIDFn func(ctx context.Context, id string) (*musichub.Album, error)
}

func (s *AlbumService) FindMeAlbums(ctx context.Context) ([]*musichub.Album, error) {
	return s.FindMeAlbumsFn(ctx)
}

func (s *AlbumService) FindAlbumByID(ctx context.Context, id string) (*musichub.Album, error) {
	return s.FindMeAlbumByIDFn(ctx, id)
}
