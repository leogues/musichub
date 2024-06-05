import { ArtistInfo } from "app/artist/artist";

import { WritableSignal } from "@angular/core";
import {
    DateContent, IndexContent, TextContent, TitleContent
} from "@components/music-table/music-table";
import { SupportedSources } from "@type/providerAuth";
import { Track, TrackResponse } from "@type/track";

export type Album = {
  isSelected: WritableSignal<boolean>;
  id: string;
  type: string;
  platform: SupportedSources;
  title: string;
  artist: ArtistInfo;
  link: string;
  picture: string;
  release_date: string;
  total_tracks: number;
};

export type AlbumWithTracks = Omit<Album, "isSelected"> & {
  tracks: Track[];
};

export type AlbumsResponse = Omit<Album, "isSelected">[];

export type AlbumResponseWithTracks = Omit<AlbumWithTracks, "tracks"> & {
  tracks: TrackResponse[];
};

export type ProvidersAlbums = Record<SupportedSources, Album[]>;

export type AlbumsTableData = {
  id: string;
  isSelected: WritableSignal<boolean>;
  link: string;
  content: [
    IndexContent,
    TitleContent,
    TextContent,
    TextContent,
    TextContent,
    DateContent,
  ];
};
