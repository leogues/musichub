package musichub

import "context"

type Album struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Platform    string      `json:"platform"`
	Title       string      `json:"title"`
	Artist      *ArtistInfo `json:"artist,omitempty"`
	Link        string      `json:"link"`
	Picture     string      `json:"picture"`
	ReleaseDate string      `json:"release_date"`
	TotalTracks int         `json:"total_tracks"`
	Tracks      []*Track    `json:"tracks,omitempty"`
}

// AlbumService provides services for albums
type AlbumService interface {
	// FindMeAlbums returns the albums for the current user
	FindMeAlbums(ctx context.Context) ([]*Album, error)
	// FindAlbumByID returns the album for the given id
	FindAlbumByID(ctx context.Context, id string) (*Album, error)
}

type AlbumInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Link  string `json:"link"`
}
