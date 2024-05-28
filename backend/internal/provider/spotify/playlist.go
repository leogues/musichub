package spotify

type Playlist struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	Name         string       `json:"name"`
	Images       []Image      `json:"images"`
	Description  string       `json:"description"`
	Public       bool         `json:"public"`
	ExternalURLs ExternalURLs `json:"external_urls"`
	Owner        struct {
		DisplayName  string       `json:"display_name"`
		ExternalURLs ExternalURLs `json:"external_urls"`
	}
	Tracks struct {
		Total int `json:"total"`
	}
}

type MePlaylistsResponse struct {
	Items []Playlist `json:"items"`
}

type PlaylistResponse struct {
	Playlist
	Tracks struct {
		Total int `json:"total"`
		Items []struct {
			Track `json:"track"`
		} `json:"items"`
	}
}
