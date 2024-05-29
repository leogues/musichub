import { PlaylistService } from "app/playlist/playlist.service";
import { PlaylistDetailComponent } from "app/playlist/ui/playlist-detail/playlist-detail.component";
import { PlaylistTableComponent } from "app/playlist/ui/playlist-table/playlist-table.component";

import { Component, computed, inject, input, signal } from "@angular/core";
import { DataInfoComponent } from "@components/data-info/data-info.component";
import { LayoutComponent } from "@components/layout/layout.component";
import { MusicPlayerComponent } from "@components/music-player/music-player.component";
import { SelectionFilterToolbarComponent } from "@components/selection-filter-toolbar/selection-filter-toolbar.component";

@Component({
  selector: "app-playlist",
  standalone: true,
  imports: [
    MusicPlayerComponent,
    LayoutComponent,
    PlaylistTableComponent,
    SelectionFilterToolbarComponent,
    DataInfoComponent,
    PlaylistDetailComponent,
  ],
  templateUrl: "./playlist.component.html",
})
export class PlaylistComponent {
  private playlistService = inject(PlaylistService);
  platform = input("", { alias: "provider" });
  playlistId = input("", { alias: "id" });

  textFilter = signal("");
  optionFilter = signal("");

  protected playlist = this.playlistService.playlist.data;
  protected playlistIsFetching = this.playlistService.playlist.isFetching;

  protected playlistTracks = computed(() => this.playlist()?.tracks || []);
  protected playlistTracksCount = computed(() => this.playlistTracks().length);
  protected selectedPlaylistTracks = computed(() => {
    return this.playlistTracks().filter((track) => track.isSelected());
  });

  protected playlistTracksFiltered = computed(() => {
    const text = this.textFilter().toLowerCase();

    if (!text) return this.playlistTracks();
    return this.playlistTracks().filter((track) => {
      return (
        track.title.toLowerCase().includes(text.toLowerCase()) ||
        track.artist.name.toLowerCase().includes(text.toLowerCase()) ||
        track.album.title.toLowerCase().includes(text.toLowerCase())
      );
    });
  });
  protected playlistTracksFilteredCount = computed(() => {
    return this.playlistTracksFiltered().length;
  });
  protected playlistTracksFilteredSelectedCount = computed(() => {
    return this.playlistTracksFiltered().reduce((acc, track) => {
      return track.isSelected() ? acc + 1 : acc;
    }, 0);
  });

  protected isAllSelected = computed(() => {
    if (this.playlistTracksFilteredCount() === 0) return false;

    return (
      this.playlistTracksFilteredCount() ===
      this.playlistTracksFilteredSelectedCount()
    );
  });

  toggleSelectAllFiltered() {
    const isAllSelected = this.isAllSelected();
    const playlistTracksFiltered = this.playlistTracksFiltered();

    playlistTracksFiltered.forEach((track) => {
      track.isSelected.set(!isAllSelected);
    });
  }

  handleRefreshData() {
    this.playlistService.fetchPlaylist(this.platform(), this.playlistId());
  }

  ngOnInit() {
    this.playlistService.fetchPlaylist(this.platform(), this.playlistId());
  }

  ngOnDestroy() {
    this.playlistService.destroy();
  }
}
