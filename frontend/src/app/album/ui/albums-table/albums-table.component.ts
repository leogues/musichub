import { Album } from "app/album/album";

import { Component, input } from "@angular/core";
import { THeader } from "@components/music-table/music-table";
import { MusicTableComponent } from "@components/music-table/music-table.component";

import { albumsToTableData } from "./toTableData";

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
    label: "Artista",
    positionCenter: false,
    labelType: "text",
    canOrder: true,
  },

  {
    contentType: "text",
    label: "Faixas",
    positionCenter: false,
    labelType: "text",
    canOrder: true,
  },
  {
    contentType: "date",
    label: "Lançamento",
    positionCenter: false,
    labelType: "text",
    canOrder: true,
  },
];

@Component({
  selector: "app-albums-table",
  standalone: true,
  imports: [MusicTableComponent],
  templateUrl: "./albums-table.component.html",
})
export class AlbumsTableComponent {
  meAlbumsTableFormated = input.required({
    alias: "albums",
    transform: (albums: Album[]) => {
      return albumsToTableData(albums);
    },
  });
  protected meAlbumsTableHeader = tableHeader;
  isLoading = input(false);
}
