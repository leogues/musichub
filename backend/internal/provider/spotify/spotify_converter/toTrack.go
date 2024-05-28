package spotifyconverter

import (
	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/provider/spotify"
)

func toTrack(spotifyTrack spotify.Track, spotifyAlbum *spotify.Album) *musichub.Track {
	if spotifyAlbum != nil {
		spotifyTrack.Album = spotifyAlbum
	}

	return &musichub.Track{
		ID:          spotifyTrack.ID,
		Type:        spotifyTrack.Type,
		Platform:    spotify.Platform,
		Title:       spotifyTrack.Name,
		Artist:      toArtistInfo(spotifyTrack.Artists[0]),
		Link:        spotifyTrack.ExternalURLs.Spotify,
		Picture:     spotifyTrack.Album.Images[0].URL,
		Album:       toAlbumInfo(spotifyTrack.Album),
		DurationMS:  spotifyTrack.DurationMS,
		Preview:     spotifyTrack.PreviewURL,
		ReleaseDate: spotifyTrack.Album.ReleaseDate,
	}
}

func MeTracksResponse(response spotify.MeTracksResponse) []*musichub.Track {
	tracks := make([]*musichub.Track, 0, len(response.Items))

	for _, spotifyTrack := range response.Items {
		track := toTrack(spotifyTrack.Track, nil)
		tracks = append(tracks, track)
	}

	return tracks
}
