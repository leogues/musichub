package gin_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/leogues/MusicSyncHub/internal/http/gin/session"
)

func TestGinServer_oAuthGoogle(t *testing.T) {
	s := MustOpenServer(t)
	defer MustCloseServer(t, s)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Get(BaseURL + "/auth/google")
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
		t.Fatal("expected oauth state in session")
	}

	if loc, err := url.Parse(res.Header.Get("Location")); err != nil {
		t.Fatal(err)
	} else if got, want := loc.Host, `accounts.google.com`; got != want {
		t.Fatalf("Location.Host=%v, want %v", got, want)
	} else if got, want := loc.Path, `/o/oauth2/auth`; got != want {
		t.Fatalf("Location.Path=%v, want %v", got, want)
	} else if got, want := loc.Query().Get("client_id"), TestGoogleClientID; got != want {
		t.Fatalf("Location.Query.client_id=%v, want %v", got, want)
	} else if got, want := loc.Query().Get("state"), session.State; got != want {
		t.Fatalf("Location.Query.state=%v, want %v", got, want)
	}

}
