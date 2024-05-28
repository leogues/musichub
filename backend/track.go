package musichub

import "context"

type Track struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Platform    string      `json:"platform"`
	Title       string      `json:"title"`
	Artist      *ArtistInfo `json:"artist,omitempty"`
	Album       *AlbumInfo  `json:"album,omitempty"`
	DurationMS  int         `json:"duration_ms"`
	Link        string      `json:"link"`
	Preview     string      `json:"preview"`
	Picture     string      `json:"picture"`
	ReleaseDate string      `json:"release_date"`
}

// TrackService provides services for tracks
type TrackService interface {
	// FindMeTracks returns the tracks for the current user
	FindMeTracks(ctx context.Context) ([]*Track, error)
}
