import { ProviderAuthResponse } from './providerAuth';

export type User = {
  id: number;
  name: string;
  email: string;
  auths: Auth[];
  provider_auths: ProviderAuthResponse[];
};

export type Auth = {
  id: number;
  source: string;
  source_id: string;
};
