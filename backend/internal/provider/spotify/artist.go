package spotify

type Artist struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	Name         string       `json:"name"`
	Images       []Image      `json:"images"`
	ExternalURLs ExternalURLs `json:"external_urls"`
	Followers    Followers    `json:"followers"`
}

type MeArtistsResponse struct {
	Artists struct {
		Total int      `json:"total"`
		Items []Artist `json:"items"`
	} `json:"artists"`
}
