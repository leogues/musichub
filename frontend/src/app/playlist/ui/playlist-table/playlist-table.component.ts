import { tracksToTableData } from "app/track/ui/track-table/toTableData";

import { Component, input } from "@angular/core";
import { THeader } from "@components/music-table/music-table";
import { MusicTableComponent } from "@components/music-table/music-table.component";
import { Track, TracksTableData } from "@type/track";

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
  selector: "app-playlist-table",
  standalone: true,
  imports: [MusicTableComponent],
  templateUrl: "./playlist-table.component.html",
})
export class PlaylistTableComponent {
  isLoading = input<boolean>(false);
  playlistTableHeader = tableHeader;
  tableTracks = input.required<TracksTableData[], Track[]>({
    alias: "playlistTracks",
    transform: tracksToTableData,
  });
}
