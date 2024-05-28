import { TrackService } from "app/track/track.service";
import { TrackTableComponent } from "app/track/ui/track-table/track-table.component";

import { Component, computed, inject, input, signal } from "@angular/core";
import { Router } from "@angular/router";
import { DataInfoComponent } from "@components/data-info/data-info.component";
import { LayoutComponent } from "@components/layout/layout.component";
import { MusicPlayerComponent } from "@components/music-player/music-player.component";
import { SelectionFilterToolbarComponent } from "@components/selection-filter-toolbar/selection-filter-toolbar.component";

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
  protected selectedMeTracks = computed(() =>
    this.meTracksData().filter((track) => track.isSelected()),
  );

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
  protected meTracksFilteredCount = computed(
    () => this.meTracksFiltered().length,
  );
  protected meTracksFilteredSelectedCount = computed(() =>
    this.meTracksFiltered().reduce(
      (acc, track) => (track.isSelected() ? acc + 1 : acc),
      0,
    ),
  );

  protected isAllSelected = computed(() => {
    if (this.meTracksFilteredCount() === 0) return false;

    return (
      this.meTracksFilteredCount() === this.meTracksFilteredSelectedCount()
    );
  });

  toggleSelectAllFiltered() {
    const isAllSelected = this.isAllSelected();
    const meTracksFiltered = this.meTracksFiltered();

    meTracksFiltered.forEach((track) => track.isSelected.set(!isAllSelected));
  }

  handleRefreshData() {
    this.trackService.refreshTracks();
  }
}
