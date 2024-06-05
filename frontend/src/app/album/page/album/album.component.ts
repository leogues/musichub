import { AlbumService } from "app/album/album.service";
import { AlbumDetailComponent } from "app/album/ui/album-detail/album-detail.component";
import { AlbumTableComponent } from "app/album/ui/album-table/album-table.component";

import { Component, computed, inject, input, signal } from "@angular/core";
import { DataInfoComponent } from "@components/data-info/data-info.component";
import { LayoutComponent } from "@components/layout/layout.component";
import { MusicPlayerComponent } from "@components/music-player/music-player.component";
import { MusicTableService } from "@components/music-table/music-table.service";
import {
    SelectionFilterToolbarComponent
} from "@components/selection-filter-toolbar/selection-filter-toolbar.component";

@Component({
  selector: "app-album",
  standalone: true,
  imports: [
    LayoutComponent,
    MusicPlayerComponent,
    DataInfoComponent,
    SelectionFilterToolbarComponent,
    AlbumTableComponent,
    AlbumDetailComponent,
  ],
  templateUrl: "./album.component.html",
})
export class AlbumComponent {
  private albumService = inject(AlbumService);
  private tableMusicService = inject(MusicTableService);
  platform = input("", { alias: "provider" });
  albumId = input("", { alias: "id" });

  protected textFilter = signal("");
  protected optionFilter = signal("");

  protected album = this.albumService.album;
  protected albumData = this.album.data;

  protected albumTracks = computed(() => this.albumData()?.tracks || []);
  protected albumTracksCount = computed(() => this.albumTracks().length);

  protected albumTracksFiltered = computed(() => {
    const text = this.textFilter().toLowerCase();

    if (!text) return this.albumTracks();
    return this.albumTracks().filter((track) => {
      return (
        track.title.toLowerCase().includes(text.toLowerCase()) ||
        track.artist.name.toLowerCase().includes(text.toLowerCase())
      );
    });
  });
  protected albumTracksFilteredCount = this.tableMusicService.dataCount;
  protected albumTracksFilteredSelectedCount =
    this.tableMusicService.dataSelectedCount;
  protected isAllSelected = this.tableMusicService.isAllSelected;

  toggleSelectAllFiltered() {
    this.tableMusicService.toggleSelectAll();
  }

  handleRefreshData() {
    this.albumService.fetchAlbum(this.platform(), this.albumId());
  }

  ngOnInit() {
    this.albumService.fetchAlbum(this.platform(), this.albumId());
  }

  ngOnDestroy() {
    this.albumService.destroy();
  }
}
