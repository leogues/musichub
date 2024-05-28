package gin_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	musichub "github.com/leogues/MusicSyncHub"
)

func TestGinServer_getMe(t *testing.T) {
	s := MustOpenServer(t)
	defer MustCloseServer(t, s)

	user0 := &musichub.User{
		ID:   1,
		Name: "jhon",
	}

	ctx0 := musichub.NewContextWithUser(context.Background(), user0)

	s.UserService.FindUsersByIDFn = func(ctx context.Context, id int) (*musichub.User, error) {
		return user0, nil
	}

	resp, err := http.DefaultClient.Do(s.MustNewRequest(t, ctx0, http.MethodGet, "/me", nil))
	if err != nil {
		t.Fatal(err)
	} else if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("StatusCode=%v, want %v", got, want)
	}

	user, err := parseResponse[musichub.User](t, resp)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&user, user0) {
		t.Fatalf("response=%#v, want %#v", user, user0)
	}

}
