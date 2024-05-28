package spotify

type Track struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	Name         string       `json:"name"`
	Artists      []*Artist    `json:"artists"`
	Album        *Album       `json:"album"`
	DurationMS   int          `json:"duration_ms"`
	ExternalURLs ExternalURLs `json:"external_urls"`
	PreviewURL   string       `json:"preview_url"`
}

type MeTracksResponse struct {
	Items []struct {
		Track Track `json:"track"`
	}
}
