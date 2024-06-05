import { CommonModule } from "@angular/common";
import { Component, input, output, signal } from "@angular/core";

import { THeader, THeaderWithOrder, TOrderEnum } from "../music-table";

@Component({
  selector: "app-music-table-header",
  standalone: true,
  imports: [CommonModule],
  templateUrl: "./music-table-header.component.html",
})
export class MusicTableHeaderComponent {
  header = input.required<THeaderWithOrder[], THeader[]>({
    transform: (header) => {
      return header.map((headerItem, index) => {
        return {
          ...headerItem,
          index,
          order: null,
        };
      });
    },
  });
  protected sortedHeader = signal<THeaderWithOrder | null>(null);
  sort = output<THeaderWithOrder>();

  handleSort(header: THeaderWithOrder) {
    if (!header.canOrder) {
      return;
    }

    if (!header.order || this.sortedHeader()?.index !== header.index) {
      header.order = TOrderEnum.ASC;
    } else if (header.order === TOrderEnum.ASC) {
      header.order = TOrderEnum.DESC;
    } else {
      header.order = null;
    }

    this.sortedHeader.set(header);
    this.sort.emit(header);
  }
}
