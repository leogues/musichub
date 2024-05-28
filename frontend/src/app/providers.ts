import { ProviderAuthInfo, SupportedSources } from '@type/providerAuth';

export const suportedProviders: Record<SupportedSources, ProviderAuthInfo> = {
  spotify: {
    source: "spotify",
    picture: "assets/images/spotify.png",
    label: "Spotify",
  },
  youtube: {
    source: "youtube",
    picture: "assets/images/youtube.png",
    label: "YouTube",
  },
};
