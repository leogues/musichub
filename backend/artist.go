package musichub

import "context"

type Artist struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Platform string `json:"platform"`
	Name     string `json:"name"`
	Fans     int    `json:"fans"`
	Link     string `json:"link"`
	Picture  string `json:"picture"`
}

// ArtistService provides services for artists
type ArtistService interface {
	// FindMeArtists returns the artists for the current user
	FindMeArtists(ctx context.Context) ([]*Artist, error)
}

type ArtistInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}
