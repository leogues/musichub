import { ScrollingModule } from "@angular/cdk/scrolling";
import { CommonModule } from "@angular/common";
import { Component, input } from "@angular/core";
import { OrderColumnPipe } from "@pipe/order-column.pipe";

import { TData, THeaderWithOrder } from "../music-table";
import { MusicTableRowComponent } from "../music-table-row/music-table-row.component";

@Component({
  selector: "app-music-table-body",
  standalone: true,
  imports: [
    CommonModule,
    ScrollingModule,
    MusicTableRowComponent,
    OrderColumnPipe,
  ],
  templateUrl: "./music-table-body.component.html",
})
export class MusicTableBodyComponent {
  data = input.required<TData[]>();
  isLoading = input<boolean>(false);
  columnWillOrdered = input.required<THeaderWithOrder | null>();

  trackByFn(_: number, item: TData) {
    return item.id;
  }
}
