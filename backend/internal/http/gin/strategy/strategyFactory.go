package strategy

import (
	"context"

	musichub "github.com/leogues/MusicSyncHub"
	"golang.org/x/oauth2"
)

type OAuthStrategyFactory struct {
	strategys *OAuthStrategys
}

func NewOAuthStrategyFactory(strategys *OAuthStrategys) *OAuthStrategyFactory {
	return &OAuthStrategyFactory{
		strategys: strategys,
	}
}

func (f *OAuthStrategyFactory) Strategy(ctx context.Context, providerName string) *oauth2.Config {
	switch providerName {
	case musichub.AuthSourceSpotify:
		return f.strategys.Spotify
	default:
		return nil
	}
}
