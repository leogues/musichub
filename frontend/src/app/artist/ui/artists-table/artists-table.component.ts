import { Artist } from 'app/artist/artist';

import { Component, input } from '@angular/core';
import { THeader } from '@components/music-table/music-table';
import { MusicTableComponent } from '@components/music-table/music-table.component';

import { artistsToTableData } from './toTableData';

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
    label: "Fãs",
    positionCenter: false,
    labelType: "text",
    canOrder: true,
  },
];

@Component({
  selector: "app-artists-table",
  standalone: true,
  imports: [MusicTableComponent],
  templateUrl: "./artists-table.component.html",
})
export class ArtistsTableComponent {
  meArtistsTableFormated = input.required({
    alias: "artists",
    transform: (artists: Artist[]) => {
      return artistsToTableData(artists);
    },
  });
  protected meArtistsTableHeader = tableHeader;
  isLoading = input(false);
}
