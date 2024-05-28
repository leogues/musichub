import { Component, input } from '@angular/core';

import { TData } from '../music-table';
import { LinkContainerComponent } from './link-container/link-container.component';
import { TableDateComponent } from './table-date/table-date.component';
import { TableDurationComponent } from './table-duration/table-duration.component';
import { TableIndexComponent } from './table-index/table-index.component';
import { TableSelectComponent } from './table-select/table-select.component';
import { TableTextComponent } from './table-text/table-text.component';
import { TableTitleComponent } from './table-title/table-title.component';
import { TableTrackComponent } from './table-track/table-track.component';

@Component({
  selector: "app-music-table-row",
  standalone: true,
  imports: [
    TableSelectComponent,
    TableIndexComponent,
    TableTitleComponent,
    TableTextComponent,
    TableTrackComponent,
    TableDurationComponent,
    TableDateComponent,
    LinkContainerComponent,
  ],
  templateUrl: "./music-table-row.component.html",
})
export class MusicTableRowComponent {
  data = input.required<TData>();
}
