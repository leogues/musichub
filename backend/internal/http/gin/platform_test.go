package gin_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	musichub "github.com/leogues/MusicSyncHub"
)

func PreConfigPlatform(t *testing.T) (*Server, context.Context) {
	t.Helper()
	s := MustOpenServer(t)

	user0 := &musichub.User{
		ID:   1,
		Name: "jhon",
	}

	ctx0 := musichub.NewContextWithUser(context.Background(), user0)

	s.UserService.FindUsersByIDFn = func(ctx context.Context, id int) (*musichub.User, error) {
		return user0, nil
	}

	return s, ctx0

}

func TestGinServer_getPlaforms(t *testing.T) {
	s, ctx0 := PreConfigPlatform(t)
	defer MustCloseServer(t, s)

	platform0 := &musichub.Platform{
		Id:   "1",
		Name: "platform",
	}

	s.ProviderService.PlatformsFn = func() []*musichub.Platform {
		return []*musichub.Platform{platform0}
	}

	resp, err := http.DefaultClient.Do(s.MustNewRequest(t, ctx0, http.MethodGet, "/platforms/list", nil))
	if err != nil {
		t.Fatal(err)
	} else if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("StatusCode=%v, want %v", got, want)
	}

	platforms, err := parseResponse[[]musichub.Platform](t, resp)
	if err != nil {
		t.Fatal(err)
	} else if len(platforms) != 1 {
		t.Fatal("expected platforms")
	} else if !reflect.DeepEqual(&platforms[0], platform0) {
		t.Fatalf("response=%#v, want %#v", &platforms[0], platform0)
	}

}
