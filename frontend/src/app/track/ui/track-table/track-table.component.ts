import { Component, input } from "@angular/core";
import { THeader } from "@components/music-table/music-table";
import { MusicTableComponent } from "@components/music-table/music-table.component";
import { Track, TracksTableData } from "@type/track";

import { tracksToTableData } from "./toTableData";

const tableHeader: THeader[] = [
  {
    contentType: "index",
    label: "#",
    isHidden: true,
    positionCenter: false,
    labelType: "text",
  },
  {
    contentType: "title",
    label: "Título",
    positionCenter: false,
    labelType: "text",
    canOrder: true,
  },
  {
    contentType: "text",
    label: "Serviço",
    positionCenter: false,
    labelType: "text",
    canOrder: true,
  },
  {
    contentType: "text",
    label: "Álbum",
    positionCenter: false,
    labelType: "text",
    canOrder: true,
  },

  {
    contentType: "track",
    label: "Preview",
    positionCenter: true,
    labelType: "text",
  },
  {
    contentType: "duration",
    positionCenter: true,
    label: "assets/images/timer.png",
    labelType: "image",
    canOrder: true,
  },
];

@Component({
  selector: "app-track-table",
  standalone: true,
  imports: [MusicTableComponent],
  templateUrl: "./track-table.component.html",
})
export class TrackTableComponent {
  isLoading = input<boolean>(false);
  trackTableHeader = tableHeader;
  tableTracks = input.required<TracksTableData[], Track[]>({
    alias: "tracks",
    transform: tracksToTableData,
  });
}
