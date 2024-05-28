package spotify

type ExternalURLs struct {
	Spotify string `json:"spotify"`
}

type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type Followers struct {
	Href  string `json:"href"`
	Total int    `json:"total"`
}
