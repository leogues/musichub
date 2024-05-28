package spotifyconverter

import (
	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/provider/spotify"
)

func toArtist(spotifyArtist spotify.Artist) *musichub.Artist {
	return &musichub.Artist{
		ID:       spotifyArtist.ID,
		Type:     spotifyArtist.Type,
		Platform: spotify.Platform,
		Name:     spotifyArtist.Name,
		Fans:     spotifyArtist.Followers.Total,
		Link:     spotifyArtist.ExternalURLs.Spotify,
		Picture:  spotifyArtist.Images[0].URL,
	}
}

func toArtistInfo(spotifyArtist *spotify.Artist) *musichub.ArtistInfo {
	if spotifyArtist == nil {
		return nil
	}
	return &musichub.ArtistInfo{
		ID:   spotifyArtist.ID,
		Name: spotifyArtist.Name,
		Link: spotifyArtist.ExternalURLs.Spotify,
	}
}

func ArtistsResponse(response spotify.MeArtistsResponse) []*musichub.Artist {
	artists := make([]*musichub.Artist, 0, len(response.Artists.Items))

	for _, spotifyArtist := range response.Artists.Items {
		artist := toArtist(spotifyArtist)
		artists = append(artists, artist)
	}

	return artists
}
