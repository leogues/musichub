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

func PreConfigTrack(t *testing.T) (*Server, context.Context) {
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

func TestGinServer_getMeTracks(t *testing.T) {
	s, ctx0 := PreConfigTrack(t)
	defer MustCloseServer(t, s)

	track0 := &musichub.Track{
		ID:    "1",
		Title: "title",
	}

	TrackService := &servicemock.TrackService{}

	TrackService.FindMeTracksFn = func(ctx context.Context) ([]*musichub.Track, error) {
		return []*musichub.Track{track0}, nil
	}
	s.TrackProviderFactory.TrackService = TrackService

	resp, err := http.DefaultClient.Do(s.MustNewRequest(t, ctx0, http.MethodGet, "/me/platforms/spotify/track", nil))
	if err != nil {
		t.Fatal(err)
	} else if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("StatusCode=%v, want %v", got, want)
	}

	tracks, err := parseResponse[[]musichub.Track](t, resp)
	if err != nil {
		t.Fatal(err)
	} else if len(tracks) != 1 {
		t.Fatalf("response=%#v, want %#v", tracks, []*musichub.Track{track0})
	} else if !reflect.DeepEqual(&tracks[0], track0) {
		t.Fatalf("response=%#v, want %#v", &tracks[0], track0)
	}

}
