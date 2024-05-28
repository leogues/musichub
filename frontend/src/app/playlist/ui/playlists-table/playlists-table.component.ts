import { Component, Input, input } from '@angular/core';
import { THeader } from '@components/music-table/music-table';
import { MusicTableComponent } from '@components/music-table/music-table.component';

import { Playlist, PlaylistsTableData } from '../../playlist';
import { playlistsToTableData } from './toTableData';

const tableHeader: THeader[] = [
  {
    contentType: "select",
    label: "checkbox",
    labelType: "text",
    positionCenter: true,
    isHidden: true,
  },
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
    label: "Faixas",
    positionCenter: false,
    labelType: "text",
    canOrder: true,
  },
  {
    contentType: "text",
    label: "Criador",
    positionCenter: false,
    labelType: "text",
    canOrder: true,
  },
  {
    contentType: "text",
    label: "Tipos",
    positionCenter: false,
    labelType: "text",
    canOrder: true,
  },
];

@Component({
  selector: "app-playlists-table",
  standalone: true,
  imports: [MusicTableComponent],
  templateUrl: "./playlists-table.component.html",
})
export class PlaylistsTableComponent {
  @Input({
    required: true,
    alias: "playlists",
    transform: (playlists: Playlist[]) => {
      return playlistsToTableData(playlists);
    },
  })
  mePlaylistsTableFormated!: PlaylistsTableData[];
  protected mePlaylistsTableHeader = tableHeader;
  isLoading = input(false);
}
