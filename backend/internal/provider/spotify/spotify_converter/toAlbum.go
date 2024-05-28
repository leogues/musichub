package spotifyconverter

import (
	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/provider/spotify"
)

func toAlbum(a spotify.Album) *musichub.Album {
	return &musichub.Album{
		ID:          a.ID,
		Type:        a.Type,
		Platform:    spotify.Platform,
		Title:       a.Name,
		Artist:      toArtistInfo(a.Artists[0]),
		Link:        a.ExternalURLs.Spotify,
		Picture:     a.Images[0].URL,
		ReleaseDate: a.ReleaseDate,
		TotalTracks: a.Tracks.Total,
	}
}

func toAlbumInfo(a *spotify.Album) *musichub.AlbumInfo {
	if a == nil {
		return nil
	}
	return &musichub.AlbumInfo{
		ID:    a.ID,
		Title: a.Name,
		Link:  a.ExternalURLs.Spotify,
	}
}

func MeAlbumsResponse(response spotify.MeAlbumsResponse) []*musichub.Album {
	albums := make([]*musichub.Album, 0, len(response.Items))

	for _, spotifyAlbum := range response.Items {
		album := toAlbum(spotifyAlbum.Album)
		albums = append(albums, album)
	}

	return albums
}

func AlbumResponse(response spotify.AlbumResponse) *musichub.Album {
	response.Album.Tracks.Total = response.TotalTracks
	album := toAlbum(response.Album)
	album.Tracks = make([]*musichub.Track, 0, len(response.Tracks.Items))

	for _, spotifyTrack := range response.Tracks.Items {
		album.Tracks = append(album.Tracks, toTrack(spotifyTrack, &response.Album))
	}

	return album
}
