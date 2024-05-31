import { AlbumsService } from "app/album/albums.service";
import { AlbumsTableComponent } from "app/album/ui/albums-table/albums-table.component";

import { Component, computed, inject, input, signal } from "@angular/core";
import { Router } from "@angular/router";
import { DataInfoComponent } from "@components/data-info/data-info.component";
import { LayoutComponent } from "@components/layout/layout.component";
import { SelectionFilterToolbarComponent } from "@components/selection-filter-toolbar/selection-filter-toolbar.component";

@Component({
  selector: "app-albums",
  standalone: true,
  imports: [
    LayoutComponent,
    DataInfoComponent,
    SelectionFilterToolbarComponent,
    AlbumsTableComponent,
  ],
  templateUrl: "./albums.component.html",
})
export class AlbumsComponent {
  private albumsService = inject(AlbumsService);
  platform = input("", { alias: "provider" });
  router = inject(Router);

  protected meAlbums = this.albumsService.meAlbums;

  protected textFilter = signal("");
  protected optionFilter = signal<string[]>([]);

  protected meAlbumsSources = computed(() => {
    if (this.meAlbums.isPending()) return [];

    const availablePlatforms = Object.keys(this.meAlbums.data());

    if (!this.platform()) return availablePlatforms;

    const hasDataForRoutePlatform = availablePlatforms.includes(
      this.platform(),
    );

    if (!hasDataForRoutePlatform) {
      this.router.navigate(["./album"]);
      return [];
    }

    return [this.platform()];
  });

  protected meAlbumsData = computed(() =>
    Object.values(this.meAlbums.data()).flat(),
  );
  protected meAlbumsCount = computed(() => this.meAlbumsData().length);
  protected selectedMeAlbums = computed(() =>
    this.meAlbumsData().filter((album) => album.isSelected()),
  );
  protected selectedMeAlbumsCount = computed(
    () => this.selectedMeAlbums().length,
  );

  protected meAlbumsFiltered = computed(() => {
    const text = this.textFilter().toLowerCase();
    const optionsFilter = this.optionFilter();
    const optionsToFilter = this.meAlbumsSources();

    const isTextEmpty = !text;
    const areAllOptionsSelected =
      optionsFilter.length === optionsToFilter.length;

    if (!isTextEmpty && areAllOptionsSelected) return this.meAlbumsData();

    return this.meAlbumsData().filter((album) => {
      const isPlatformIncluded = optionsFilter.includes(album.platform);
      const isTextMatch =
        album.title.toLowerCase().includes(text) ||
        album.artist.name.toLowerCase().includes(text);
      return isPlatformIncluded && isTextMatch;
    });
  });
  protected meAlbumsFilteredCount = computed(
    () => this.meAlbumsFiltered().length,
  );
  protected meAlbumsFilteredSelectedCount = computed(() =>
    this.meAlbumsFiltered().reduce(
      (acc, album) => (album.isSelected() ? acc + 1 : acc),
      0,
    ),
  );

  protected isAllSelected = computed(() => {
    if (this.meAlbumsFilteredCount() === 0) return false;

    return (
      this.meAlbumsFilteredCount() === this.meAlbumsFilteredSelectedCount()
    );
  });

  toggleSelectAllFiltered() {
    const isAllSelected = this.isAllSelected();
    const meAlbumsFiltered = this.meAlbumsFiltered();

    meAlbumsFiltered.forEach((album) => {
      album.isSelected.set(!isAllSelected);
    });
  }

  handleRefreshData() {
    this.albumsService.refreshAlbums();
  }
}
