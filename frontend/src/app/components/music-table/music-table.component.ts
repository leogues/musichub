import { Component, computed, input, signal, WritableSignal } from '@angular/core';

import { TData, THeader, THeaderWithOrder } from './music-table';
import { MusicTableBodyComponent } from './music-table-body/music-table-body.component';
import { MusicTableHeaderComponent } from './music-table-header/music-table-header.component';

@Component({
  selector: "app-music-table",
  standalone: true,
  imports: [MusicTableHeaderComponent, MusicTableBodyComponent],
  templateUrl: "./music-table.component.html",
})
export class MusicTableComponent {
  header = input.required<THeader[]>();
  data = input.required<TData[]>();
  isLoading = input<boolean>(false);
  colCount = computed(() => this.header().length);

  protected columnWillOrdered: WritableSignal<THeaderWithOrder | null> =
    signal(null);

  handleSort(columnWillOrdered: THeaderWithOrder) {
    this.columnWillOrdered.set({ ...columnWillOrdered });
  }
}
