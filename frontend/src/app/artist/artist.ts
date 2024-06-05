import { WritableSignal } from "@angular/core";
import {
  IndexContent,
  TextContent,
  TitleContent,
} from "@components/music-table/music-table";
import { SupportedSources } from "@type/providerAuth";

export type Artist = {
  isSelected: WritableSignal<boolean>;
  id: string;
  type: string;
  platform: SupportedSources;
  name: string;
  fans: number;
  link: string;
  picture: string;
};

export type ArtistInfo = Pick<Artist, "id" | "name" | "link">;

export type ArtistsResponse = Omit<Artist, "isSelected">[];

export type ProvidersArtists = Record<SupportedSources, Artist[]>;

export type ArtistsTableData = {
  id: string;
  isSelected: WritableSignal<boolean>;
  content: [IndexContent, TitleContent, TextContent, TextContent];
};
