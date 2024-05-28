import { WritableSignal } from "@angular/core";
import {
  DurationContent,
  IndexContent,
  SelectContent,
  TextContent,
  TitleContent,
  TrackContent,
} from "@components/music-table/music-table";
import { SupportedSources } from "./providerAuth";

export type Track = {
  isSelected: WritableSignal<boolean>;
  id: string;
  type: string;
  platform: SupportedSources;
  title: string;
  artist: {
    id: string;
    name: string;
    link: string;
  };
  album: {
    id: string;
    title: string;
    link: string;
  };
  duration_ms: number;
  link: string;
  preview: string | undefined;
  picture: string;
  release_date: string;
};

export type TrackResponse = Omit<Track, "isSelected">;

export type TracksResponse = TrackResponse[];

export type TracksTableData = {
  id: string;
  content: [
    SelectContent,
    IndexContent,
    TitleContent,
    TextContent,
    TextContent,
    TrackContent,
    DurationContent,
  ];
};
