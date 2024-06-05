import { computed, Injectable, signal } from "@angular/core";

import { TData } from "./music-table";

@Injectable({
  providedIn: "root",
})
export class MusicTableService {
  data = signal<TData[]>([]);

  setData(data: TData[]) {
    this.data.set(data);
  }

  dataCount = computed(() => this.data().length);
  dataSelected = computed(() =>
    this.data().filter((dataRow) => dataRow.isSelected()),
  );
  dataSelectedCount = computed(() => this.dataSelected().length);
  isAllSelected = computed(() => {
    if (this.dataCount() === 0) return false;
    return this.dataSelectedCount() === this.dataCount();
  });

  toggleSelectAll() {
    const isAllSelected = this.isAllSelected();
    this.data().forEach((dataRow) => {
      dataRow.isSelected.set(!isAllSelected);
    });
  }
}
