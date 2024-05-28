import { DatePipe } from '@angular/common';
import { Component, input } from '@angular/core';
import { DateContent } from '@components/music-table/music-table';

@Component({
  selector: "app-table-date",
  standalone: true,
  imports: [DatePipe],
  templateUrl: "./table-date.component.html",
})
export class TableDateComponent {
  data = input.required<DateContent>();
}
