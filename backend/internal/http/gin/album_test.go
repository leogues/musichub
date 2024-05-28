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

func PreConfigAlbum(t *testing.T) (*Server, context.Context) {
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

func TestGinServer_getMeAlbums(t *testing.T) {
	s, ctx0 := PreConfigAlbum(t)
	defer MustCloseServer(t, s)

	album0 := &musichub.Album{
		ID:    "1",
		Title: "name",
	}

	AlbumService := &servicemock.AlbumService{}

	AlbumService.FindMeAlbumsFn = func(ctx context.Context) ([]*musichub.Album, error) {
		return []*musichub.Album{album0}, nil
	}

	s.AlbumProviderFactory.AlbumService = AlbumService

	resp, err := http.DefaultClient.Do(s.MustNewRequest(t, ctx0, http.MethodGet, "/me/platforms/spotify/album", nil))
	if err != nil {
		t.Fatal(err)
	} else if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("StatusCode=%v, want %v", got, want)
	}

	albums, err := parseResponse[[]musichub.Album](t, resp)
	if err != nil {
		t.Fatal(err)
	} else if len(albums) != 1 {
		t.Fatalf("response=%#v, want %#v", albums, []*musichub.Album{album0})
	} else if !reflect.DeepEqual(&albums[0], album0) {
		t.Fatalf("response=%#v, want %#v", &albums[0], album0)
	}

}

func TestGinServer_getAlbum(t *testing.T) {
	s, ctx0 := PreConfigAlbum(t)
	defer MustCloseServer(t, s)

	album0 := &musichub.Album{
		ID:    "1",
		Title: "name",
	}

	AlbumService := &servicemock.AlbumService{}

	AlbumService.FindMeAlbumByIDFn = func(ctx context.Context, id string) (*musichub.Album, error) {
		return album0, nil
	}

	s.AlbumProviderFactory.AlbumService = AlbumService

	resp, err := http.DefaultClient.Do(s.MustNewRequest(t, ctx0, http.MethodGet, "/platforms/spotify/album/1", nil))
	if err != nil {
		t.Fatal(err)
	} else if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("StatusCode=%v, want %v", got, want)
	}

	album, err := parseResponse[musichub.Album](t, resp)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&album, album0) {
		t.Fatalf("response=%#v, want %#v", &album, album0)
	}
}
