package gin_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	musichub "github.com/leogues/MusicSyncHub"
	servicemock "github.com/leogues/MusicSyncHub/internal/mock/serviceMock"
)

func PreConfigPlaylist(t *testing.T) (*Server, context.Context) {
	t.Helper()
	s := MustOpenServer(t)
	expiry := time.Now().Add(time.Hour)

	providerAuth0 := &musichub.ProviderAuth{
		ID:          1,
		Source:      musichub.AuthSourceSpotify,
		AccessToken: "access_token",
		Expiry:      &expiry,
	}

	user0 := &musichub.User{
		ID:            1,
		Name:          "jhon",
		ProviderAuths: []*musichub.ProviderAuth{providerAuth0},
	}

	ctx0 := musichub.NewContextWithUser(context.Background(), user0)

	s.UserService.FindUsersByIDFn = func(ctx context.Context, id int) (*musichub.User, error) {
		return user0, nil
	}

	return s, ctx0

}

func TestGinServer_getMePlaylists(t *testing.T) {
	s, ctx0 := PreConfigPlaylist(t)
	defer MustCloseServer(t, s)

	playlist0 := &musichub.Playlist{
		ID:       "1",
		Name:     "name",
		Platform: musichub.AuthSourceSpotify,
	}

	PlaylistService := &servicemock.PlaylistService{}

	PlaylistService.FindMePlaylistsFn = func(ctx context.Context) ([]*musichub.Playlist, error) {
		return []*musichub.Playlist{playlist0}, nil
	}

	s.PlaylistProviderFactory.PlaylistService = PlaylistService

	resp, err := http.DefaultClient.Do(s.MustNewRequest(t, ctx0, http.MethodGet, "/me/platforms/spotify/playlist", nil))
	if err != nil {
		t.Fatal(err)
	} else if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("StatusCode=%v, want %v", got, want)
	}

	playlists, err := parseResponse[[]musichub.Playlist](t, resp)
	if err != nil {
		t.Fatal(err)
	} else if len(playlists) != 1 {
		t.Fatalf("response=%#v, want %#v", playlists, []*musichub.Playlist{playlist0})
	} else if !reflect.DeepEqual(&playlists[0], playlist0) {
		t.Fatalf("response=%#v, want %#v", &playlists[0], playlist0)
	}
}

func TestGinServer_getPlaylist(t *testing.T) {
	s, ctx0 := PreConfigPlaylist(t)
	defer MustCloseServer(t, s)

	playlist0 := &musichub.Playlist{
		ID:   "1",
		Name: "name",
	}

	PlaylistService := &servicemock.PlaylistService{}

	PlaylistService.FindPlaylistByIDFn = func(ctx context.Context, id string) (*musichub.Playlist, error) {
		return playlist0, nil
	}

	s.PlaylistProviderFactory.PlaylistService = PlaylistService

	resp, err := http.DefaultClient.Do(s.MustNewRequest(t, ctx0, http.MethodGet, "/platforms/spotify/playlist/1", nil))
	if err != nil {
		t.Fatal(err)
	} else if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("StatusCode=%v, want %v", got, want)
	}

	playlist, err := parseResponse[musichub.Playlist](t, resp)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&playlist, playlist0) {
		t.Fatalf("response=%#v, want %#v", &playlist, playlist0)
	}
}
