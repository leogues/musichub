import { AlbumService } from "app/album/album.service";
import { AlbumDetailComponent } from "app/album/ui/album-detail/album-detail.component";
import { AlbumTableComponent } from "app/album/ui/album-table/album-table.component";

import {
  Component,
  Input,
  computed,
  inject,
  input,
  model,
  output,
  signal,
} from "@angular/core";
import { DataInfoComponent } from "@components/data-info/data-info.component";
import { LayoutComponent } from "@components/layout/layout.component";
import { MusicPlayerComponent } from "@components/music-player/music-player.component";
import { SelectionFilterToolbarComponent } from "@components/selection-filter-toolbar/selection-filter-toolbar.component";

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
  platform = input("", { alias: "provider" });
  albumId = input("", { alias: "id" });

  textFilter = signal("");
  optionFilter = signal("");

  protected album = this.albumService.album.data;
  protected albumIsFetching = this.albumService.album.isFetching;

  protected albumTracks = computed(() => this.album()?.tracks || []);
  protected albumTracksCount = computed(() => this.albumTracks().length);
  protected selectedAlbumTracks = computed(() => {
    return this.albumTracks().filter((track) => track.isSelected());
  });

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
  protected albumTracksFilteredCount = computed(() => {
    return this.albumTracksFiltered().length;
  });
  protected albumTracksFilteredSelectedCount = computed(() => {
    return this.albumTracksFiltered().reduce((acc, track) => {
      return track.isSelected() ? acc + 1 : acc;
    }, 0);
  });

  protected isAllSelected = computed(() => {
    if (this.albumTracksFilteredCount() === 0) return false;
    return (
      this.albumTracksFilteredSelectedCount() ===
      this.albumTracksFilteredCount()
    );
  });

  toggleSelectAllFiltered() {
    const isAllSelected = this.isAllSelected();
    const playlistTracksFiltered = this.albumTracksFiltered();

    playlistTracksFiltered.forEach((track) => {
      track.isSelected.set(!isAllSelected);
    });
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
