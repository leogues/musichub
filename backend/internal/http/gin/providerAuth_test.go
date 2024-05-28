package gin_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/http/gin/session"
)

func PreConfigProviderAuth(t *testing.T) (*Server, context.Context) {
	s := MustOpenServer(t)

	user0 := &musichub.User{
		ID:   1,
		Name: "jhon",
	}

	s.UserService.FindUsersByIDFn = func(ctx context.Context, id int) (*musichub.User, error) {
		return user0, nil
	}

	ctx0 := musichub.NewContextWithUser(context.Background(), user0)

	return s, ctx0

}

func TestGinServer_oAuthSpotify(t *testing.T) {
	s, ctx0 := PreConfigProviderAuth(t)
	defer MustCloseServer(t, s)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Do(s.MustNewRequest(t, ctx0, http.MethodGet, "/auth/spotify", nil))
	if err != nil {
		t.Fatal(err)
	} else if err := res.Body.Close(); err != nil {
		t.Fatal(err)
	} else if got, want := res.StatusCode, http.StatusFound; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	cookie, _ := url.QueryUnescape(res.Cookies()[0].Value)

	session := session.NewSession()

	if err := s.sessionManager.UnmarshalSession(cookie, &session); err != nil {
		t.Fatal(err)
	} else if session.State == "" {
		t.Fatal(session)
	}

	if loc, err := url.Parse(res.Header.Get("Location")); err != nil {
		t.Fatal(err)
	} else if got, want := loc.Host, `accounts.spotify.com`; got != want {
		t.Fatalf("Location.Host=%v, want %v", got, want)
	} else if got, want := loc.Path, `/authorize`; got != want {
		t.Fatalf("Location.Path=%v, want %v", got, want)
	} else if got, want := loc.Query().Get("client_id"), TestSpotifyClientID; got != want {
		t.Fatalf("Location.Query.client_id=%v, want %v", got, want)
	} else if got, want := loc.Query().Get("state"), session.State; got != want {
		t.Fatalf("Location.Query.state=%v, want %v", got, want)
	}

}
