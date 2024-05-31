import { PlaylistsService } from "app/playlist/playlists.service";
import { PlaylistsTableComponent } from "app/playlist/ui/playlists-table/playlists-table.component";

import { Component, computed, inject, input, signal } from "@angular/core";
import { Router } from "@angular/router";
import { DataInfoComponent } from "@components/data-info/data-info.component";
import { LayoutComponent } from "@components/layout/layout.component";
import { SelectionFilterToolbarComponent } from "@components/selection-filter-toolbar/selection-filter-toolbar.component";

@Component({
  selector: "app-playlists",
  standalone: true,
  imports: [
    LayoutComponent,
    PlaylistsTableComponent,
    SelectionFilterToolbarComponent,
    DataInfoComponent,
  ],
  templateUrl: "./playlists.component.html",
})
export class PlaylistsComponent {
  private playlistsService = inject(PlaylistsService);
  platform = input("", { alias: "provider" });
  router = inject(Router);

  protected mePlaylists = this.playlistsService.mePlaylists;
  // mePlaylistsData = computed(() => this.mePlaylists.data());
  protected mePlaylistIsFetching = this.mePlaylists.isPending();
  protected textFilter = signal("");
  protected optionFilter = signal<string[]>([]);

  protected mePlaylistsSources = computed(() => {
    if (this.mePlaylists.isPending()) return [];

    const availablePlatforms = Object.keys(this.mePlaylists.data());

    if (!this.platform()) return availablePlatforms;

    const hasDataForRoutePlatform = availablePlatforms.includes(
      this.platform(),
    );

    if (!hasDataForRoutePlatform) {
      this.router.navigate(["./playlist"]);
      return [];
    }

    return [this.platform()];
  });

  protected mePlaylistsData = computed(() =>
    Object.values(this.mePlaylists.data()).flat(),
  );
  protected mePlaylistsCount = computed(() => this.mePlaylistsData().length);
  protected selectedMePlaylists = computed(() =>
    this.mePlaylistsData().filter((playlist) => playlist.isSelected()),
  );
  protected selectedMePlaylistsCount = computed(
    () => this.selectedMePlaylists().length,
  );

  protected mePlaylistsFiltered = computed(() => {
    const text = this.textFilter().toLowerCase();
    const optionsFilter = this.optionFilter();
    const optionsToFilter = this.mePlaylistsSources();

    const isTextEmpty = !text;
    const areAllOptionsSelected =
      optionsFilter.length === optionsToFilter.length;

    if (isTextEmpty && areAllOptionsSelected) return this.mePlaylistsData();

    return this.mePlaylistsData().filter((playlist) => {
      const isPlatformIncluded = optionsFilter.includes(playlist.platform);
      const isTextMatch =
        playlist.title.toLowerCase().includes(text) ||
        playlist.creator.toLowerCase().includes(text);
      return isPlatformIncluded && isTextMatch;
    });
  });
  protected mePlaylistsFilteredCount = computed(
    () => this.mePlaylistsFiltered().length,
  );
  protected mePlaylistsFilteredSelectedCount = computed(() =>
    this.mePlaylistsFiltered().reduce(
      (acc, playlist) => (playlist.isSelected() ? acc + 1 : acc),
      0,
    ),
  );

  protected isAllSelected = computed(() => {
    if (this.mePlaylistsFilteredCount() === 0) return false;

    return (
      this.mePlaylistsFilteredCount() ===
      this.mePlaylistsFilteredSelectedCount()
    );
  });

  toggleSelectAllFiltered() {
    const isAllSelected = this.isAllSelected();
    const mePlaylistsFiltered = this.mePlaylistsFiltered();

    mePlaylistsFiltered.forEach((playlist) => {
      playlist.isSelected.set(!isAllSelected);
    });
  }

  handleRefreshData() {
    this.playlistsService.refreshPlaylists();
  }
}
