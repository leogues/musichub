import { Component, input } from '@angular/core';

@Component({
  selector: "app-data-info",
  standalone: true,
  imports: [],
  templateUrl: "./data-info.component.html",
})
export class DataInfoComponent {
  total = input.required<number>();
  selectedCount = input.required<number>();
  filteredCount = input.required<number>();
}
