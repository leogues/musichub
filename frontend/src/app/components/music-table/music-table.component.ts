import {
  Component,
  computed,
  inject,
  input,
  signal,
  WritableSignal,
} from "@angular/core";
import { toObservable } from "@angular/core/rxjs-interop";

import { TData, THeader, THeaderWithOrder } from "./music-table";
import { MusicTableBodyComponent } from "./music-table-body/music-table-body.component";
import { MusicTableHeaderComponent } from "./music-table-header/music-table-header.component";
import { MusicTableService } from "./music-table.service";
import { Subject, takeUntil } from "rxjs";

@Component({
  selector: "app-music-table",
  standalone: true,
  imports: [MusicTableHeaderComponent, MusicTableBodyComponent],
  templateUrl: "./music-table.component.html",
})
export class MusicTableComponent {
  musicTableService = inject(MusicTableService);
  header = input.required<THeader[]>();
  data = input.required<TData[]>();
  dataObservable = toObservable(this.data);
  isLoading = input<boolean>(false);
  colCount = computed(() => this.header().length);

  private destroy$ = new Subject<void>();

  constructor() {
    this.dataObservable.pipe(takeUntil(this.destroy$)).subscribe((data) => {
      this.musicTableService.setData(data);
    });
  }

  protected columnWillOrdered: WritableSignal<THeaderWithOrder | null> =
    signal(null);

  handleSort(columnWillOrdered: THeaderWithOrder) {
    this.columnWillOrdered.set({ ...columnWillOrdered });
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
