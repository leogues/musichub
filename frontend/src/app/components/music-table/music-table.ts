import { WritableSignal } from '@angular/core';

export type THeader = {
  label: string;
  labelType: "text" | "image";
  positionCenter?: boolean;
  isHidden?: boolean;
  contentType: keyof ContentMapKey;
  canOrder?: boolean;
};

type TOrder = "asc" | "desc" | null;

export enum TOrderEnum {
  ASC = "asc",
  DESC = "desc",
}

export type THeaderWithOrder = THeader & {
  index: number;
  order: TOrder;
};

export type SelectContent = {
  type: "select";
  isSelected: WritableSignal<boolean>;
};

export type IndexContent = {
  type: "index";
  index: number;
};

export type TitleContent = {
  type: "title";
  imageurl: string;
  title: string;
  link?: string;
  artist?: string;
  artistlink?: string;
};

export type TextContent = {
  type: "text";
  text: string;
  class?: string;
  link?: string;
};

export type TrackContent = {
  type: "track";
  id: string;
  trackurl: string;
};

export type ImageContent = {
  type: "image";
  imageurl: string;
  link: string;
};

export type DurationContent = {
  type: "duration";
  duration: number;
};

export type DateContent = {
  type: "date";
  date: string;
};

export type TData = {
  id: string;
  link?: string;
  content: Content[];
};

export type Content =
  | SelectContent
  | IndexContent
  | TitleContent
  | TextContent
  | TrackContent
  | ImageContent
  | DurationContent
  | DateContent;

export type ContentMapKey = {
  select: SelectContent;
  index: IndexContent;
  title: TitleContent;
  text: TextContent;
  track: TrackContent;
  image: ImageContent;
  duration: DurationContent;
  date: DateContent;
};
