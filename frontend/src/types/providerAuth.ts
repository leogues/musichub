export type SupportedSources = "spotify" | "youtube";

export type ProviderAuth = {
  id?: number;
  source: SupportedSources;
  isAuthenticating?: boolean;
  color?: string;
  picture: string;
  label: string;
};

export type ProviderAuthInfo = Omit<ProviderAuth, "isAuthenticating" | "id">;

export type ProviderAuthResponse = Omit<ProviderAuth, "picture" | "label">;
