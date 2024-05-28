package spotifyconverter

import (
	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/provider/spotify"
)

func toPlaylist(spotifyPlaylist spotify.Playlist) *musichub.Playlist {
	return &musichub.Playlist{
		ID:          spotifyPlaylist.ID,
		Type:        spotifyPlaylist.Type,
		Platform:    spotify.Platform,
		Name:        spotifyPlaylist.Name,
		Description: spotifyPlaylist.Description,
		Picture:     spotifyPlaylist.Images[0].URL,
		Link:        spotifyPlaylist.ExternalURLs.Spotify,
		Creator:     spotifyPlaylist.Owner.DisplayName,
		CreatorLink: spotifyPlaylist.Owner.ExternalURLs.Spotify,
		TotalTracks: spotifyPlaylist.Tracks.Total,
		Public:      spotifyPlaylist.Public,
	}
}

func MePlaylistsResponse(response spotify.MePlaylistsResponse) []*musichub.Playlist {
	playlists := make([]*musichub.Playlist, 0, len(response.Items))

	for _, spotifyPlaylist := range response.Items {
		playlist := toPlaylist(spotifyPlaylist)
		playlists = append(playlists, playlist)
	}

	return playlists
}

func PlaylistResponse(response spotify.PlaylistResponse) *musichub.Playlist {
	response.Playlist.Tracks.Total = response.Tracks.Total
	playlist := toPlaylist(response.Playlist)
	playlist.Tracks = make([]*musichub.Track, 0, len(response.Tracks.Items))

	for _, spotifyItem := range response.Tracks.Items {
		playlist.Tracks = append(playlist.Tracks, toTrack(spotifyItem.Track, nil))
	}

	return playlist
}
