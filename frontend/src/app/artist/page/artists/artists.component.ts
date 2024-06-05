import { ArtistService } from "app/artist/artist.service";
import { ArtistsTableComponent } from "app/artist/ui/artists-table/artists-table.component";

import { Component, computed, inject, input, signal } from "@angular/core";
import { Router } from "@angular/router";
import { DataInfoComponent } from "@components/data-info/data-info.component";
import { LayoutComponent } from "@components/layout/layout.component";
import { SelectionFilterToolbarComponent } from "@components/selection-filter-toolbar/selection-filter-toolbar.component";
import { MusicTableService } from "@components/music-table/music-table.service";

@Component({
  selector: "app-artists",
  standalone: true,
  imports: [
    LayoutComponent,
    DataInfoComponent,
    SelectionFilterToolbarComponent,
    ArtistsTableComponent,
  ],
  templateUrl: "./artists.component.html",
})
export class ArtistsComponent {
  private artistService = inject(ArtistService);
  private tableMusicService = inject(MusicTableService);
  platform = input("", { alias: "provider" });
  router = inject(Router);

  protected meArtists = this.artistService.meArtists;

  protected textFilter = signal("");
  protected optionFilter = signal<string[]>([]);

  protected meArtistsSources = computed(() => {
    if (this.meArtists.isPending()) return [];

    const availablePlatforms = Object.keys(this.meArtists.data());

    if (!this.platform()) return availablePlatforms;

    const hasDataForRoutePlatform = availablePlatforms.includes(
      this.platform(),
    );

    if (!hasDataForRoutePlatform) {
      this.router.navigate(["./artist"]);
      return [];
    }

    return [this.platform()];
  });

  protected meArtistsData = computed(() =>
    Object.values(this.meArtists.data()).flat(),
  );
  protected meArtistsCount = computed(() => this.meArtistsData().length);
  protected selectedMeArtists = computed(() =>
    this.meArtistsData().filter((artist) => artist.isSelected()),
  );
  protected selectedMeArtistsCount = computed(
    () => this.selectedMeArtists().length,
  );

  protected meArtistsFiltered = computed(() => {
    const text = this.textFilter().toLowerCase();
    const optionsFilter = this.optionFilter();
    const optionsToFilter = this.meArtistsSources();

    const isTextEmpty = !text;
    const areAllOptionsSelected =
      optionsFilter.length === optionsToFilter.length;

    if (isTextEmpty && areAllOptionsSelected) return this.meArtistsData();

    return this.meArtistsData().filter((artist) => {
      const isPlatformIncluded = optionsFilter.includes(artist.platform);
      const isTextMatch = artist.name.toLowerCase().includes(text);
      return isPlatformIncluded && isTextMatch;
    });
  });
  protected meArtistsFilteredCount = this.tableMusicService.dataCount;
  protected meArtistsFilteredSelectedCount =
    this.tableMusicService.dataSelectedCount;
  protected isAllSelected = this.tableMusicService.isAllSelected;

  toggleSelectAllFiltered() {
    this.tableMusicService.toggleSelectAll();
  }

  handleRefreshData() {
    this.artistService.refreshArtists();
  }
}
