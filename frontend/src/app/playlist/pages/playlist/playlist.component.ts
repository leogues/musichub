import { PlaylistService } from "app/playlist/playlist.service";
import { PlaylistDetailComponent } from "app/playlist/ui/playlist-detail/playlist-detail.component";
import { PlaylistTableComponent } from "app/playlist/ui/playlist-table/playlist-table.component";

import {
  Component,
  computed,
  effect,
  inject,
  input,
  signal,
} from "@angular/core";
import { DataInfoComponent } from "@components/data-info/data-info.component";
import { LayoutComponent } from "@components/layout/layout.component";
import { MusicPlayerComponent } from "@components/music-player/music-player.component";
import { SelectionFilterToolbarComponent } from "@components/selection-filter-toolbar/selection-filter-toolbar.component";
import { MusicTableService } from "@components/music-table/music-table.service";

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
  private tableMusicService = inject(MusicTableService);
  platform = input("", { alias: "provider" });
  playlistId = input("", { alias: "id" });

  textFilter = signal("");
  optionFilter = signal("");

  protected playlist = this.playlistService.playlist;
  protected playlistData = this.playlist.data;

  protected playlistTracks = computed(() => this.playlistData()?.tracks || []);
  protected playlistTracksCount = computed(() => this.playlistTracks().length);

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
  protected playlistTracksFilteredCount = this.tableMusicService.dataCount;
  protected playlistTracksFilteredSelectedCount =
    this.tableMusicService.dataSelectedCount;
  protected isAllSelected = this.tableMusicService.isAllSelected;

  toggleSelectAllFiltered() {
    this.tableMusicService.toggleSelectAll();
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
