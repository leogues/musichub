package spotify

type Album struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	Name         string       `json:"name"`
	Artists      []*Artist    `json:"artists"`
	Images       []Image      `json:"images"`
	ExternalURLs ExternalURLs `json:"external_urls"`
	ReleaseDate  string       `json:"release_date"`
	Tracks       struct {
		Total int `json:"total"`
	} `json:"tracks"`
}

type MeAlbumsResponse struct {
	Items []struct {
		Album Album `json:"album"`
	}
}

type AlbumResponse struct {
	TotalTracks int `json:"total_tracks"`
	Album
	Tracks struct {
		Total int     `json:"total"`
		Items []Track `json:"items"`
	} `json:"tracks"`
}
