import { SupportedSources } from '@type/providerAuth';
import { Track } from '@type/track';

export type ProvidersTracks = Record<SupportedSources, Track[]>;
