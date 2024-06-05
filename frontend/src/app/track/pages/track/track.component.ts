import { TrackService } from "app/track/track.service";
import { TrackTableComponent } from "app/track/ui/track-table/track-table.component";

import { Component, computed, inject, input, signal } from "@angular/core";
import { Router } from "@angular/router";
import { DataInfoComponent } from "@components/data-info/data-info.component";
import { LayoutComponent } from "@components/layout/layout.component";
import { MusicPlayerComponent } from "@components/music-player/music-player.component";
import { MusicTableService } from "@components/music-table/music-table.service";
import {
    SelectionFilterToolbarComponent
} from "@components/selection-filter-toolbar/selection-filter-toolbar.component";

@Component({
  selector: "app-track",
  standalone: true,
  imports: [
    LayoutComponent,
    SelectionFilterToolbarComponent,
    DataInfoComponent,
    MusicPlayerComponent,
    TrackTableComponent,
  ],
  templateUrl: "./track.component.html",
})
export class TrackComponent {
  private trackService = inject(TrackService);
  private tableMusicService = inject(MusicTableService);
  router = inject(Router);
  platform = input("", { alias: "provider" });

  protected meTracks = this.trackService.meTracks;

  protected textFilter = signal("");
  protected optionFilter = signal<string[]>([]);

  protected meTracksSources = computed(() => {
    if (this.meTracks.isPending()) return [];

    const availablePlatforms = Object.keys(this.meTracks.data());

    if (!this.platform()) return availablePlatforms;

    const hasDataForRoutePlatform = availablePlatforms.includes(
      this.platform(),
    );

    if (!hasDataForRoutePlatform) {
      this.router.navigate(["./track"]);
      return [];
    }

    return [this.platform()];
  });

  protected meTracksData = computed(() =>
    Object.values(this.meTracks.data()).flat(),
  );
  protected meTracksCount = computed(() => this.meTracksData().length);

  protected meTracksFiltered = computed(() => {
    const text = this.textFilter().toLowerCase();
    const optionsFilter = this.optionFilter();
    const optionsToFilter = this.meTracksSources();

    const isTextEmpty = !text;
    const areAllOptionsSelected =
      optionsFilter.length === optionsToFilter.length;

    if (isTextEmpty && areAllOptionsSelected) return this.meTracksData();

    return this.meTracksData().filter((track) => {
      const isPlatformIncluded = optionsFilter.includes(track.platform);
      const isTextMatch =
        track.title.toLowerCase().includes(text) ||
        track.artist.name.toLowerCase().includes(text) ||
        track.album.title.toLowerCase().includes(text);
      return isPlatformIncluded && isTextMatch;
    });
  });
  protected meTracksFilteredCount = this.tableMusicService.dataCount;
  protected meTracksFilteredSelectedCount =
    this.tableMusicService.dataSelectedCount;
  protected isAllSelected = this.tableMusicService.isAllSelected;

  toggleSelectAllFiltered() {
    this.tableMusicService.toggleSelectAll();
  }

  handleRefreshData() {
    this.trackService.refreshTracks();
  }
}
