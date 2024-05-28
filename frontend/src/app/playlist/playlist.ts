import { WritableSignal } from '@angular/core';
import { Track, TrackResponse } from '@type/track';

import { SupportedSources } from '../../types/providerAuth';
import {
    IndexContent, SelectContent, TextContent, TitleContent
} from '../components/music-table/music-table';

export type Playlist = {
  isSelected: WritableSignal<boolean>;
  id: string;
  type: string;
  platform: SupportedSources;
  title: string;
  description: string;
  picture: string;
  link: string;
  creator: string;
  creator_link: string;
  public: boolean;
  total_tracks: number;
};

export type PlaylistWithTracks = Omit<Playlist, "isSelected"> & {
  tracks: Track[];
};

export type PlaylistsResponse = Omit<Playlist, "isSelected">[];

export type PlaylistResponseWithTracks = Omit<PlaylistWithTracks, "tracks"> & {
  tracks: TrackResponse[];
};

export type ProvidersPlaylists = Record<SupportedSources, Playlist[]>;

export type PlaylistsTableData = {
  id: string;
  link: string;
  content: [
    SelectContent,
    IndexContent,
    TitleContent,
    TextContent,
    TextContent,
    TextContent,
    TextContent,
  ];
};
