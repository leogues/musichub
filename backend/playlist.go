package musichub

import "context"

type Playlist struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Platform    string   `json:"platform"`
	Name        string   `json:"title"`
	Description string   `json:"description"`
	Picture     string   `json:"picture"`
	Link        string   `json:"link"`
	Creator     string   `json:"creator"`
	CreatorLink string   `json:"creator_link"`
	Public      bool     `json:"public"`
	TotalTracks int      `json:"total_tracks"`
	Tracks      []*Track `json:"tracks,omitempty"`
}

// PlaylistService represents a service for managing playlists.
type PlaylistService interface {
	// FindMePlaylists returns a list of playlists owned by the current user.
	FindMePlaylists(ctx context.Context) ([]*Playlist, error)

	// FindPlaylistByID returns the playlist with the specified ID.
	FindPlaylistByID(ctx context.Context, id string) (*Playlist, error)
}
