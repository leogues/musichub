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

func PreConfigArtist(t *testing.T) (*Server, context.Context) {
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

func TestGinServer_getMeArtists(t *testing.T) {
	s, ctx0 := PreConfigArtist(t)
	defer MustCloseServer(t, s)

	artist0 := &musichub.Artist{
		ID:   "1",
		Name: "artist",
	}

	ArtistService := &servicemock.ArtistService{}

	ArtistService.FindMeArtistsFn = func(ctx context.Context) ([]*musichub.Artist, error) {
		return []*musichub.Artist{artist0}, nil
	}
	s.ArtistProviderFactory.ArtistService = ArtistService

	resp, err := http.DefaultClient.Do(s.MustNewRequest(t, ctx0, http.MethodGet, "/me/platforms/spotify/artist", nil))
	if err != nil {
		t.Fatal(err)
	} else if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("StatusCode=%v, want %v", got, want)
	}

	artists, err := parseResponse[[]musichub.Artist](t, resp)
	if err != nil {
		t.Fatal(err)
	} else if len(artists) != 1 {
		t.Fatalf("response=%#v, want %#v", artists, []*musichub.Artist{artist0})
	} else if !reflect.DeepEqual(&artists[0], artist0) {
		t.Fatalf("response=%#v, want %#v", artists[0], artist0)
	}

}
