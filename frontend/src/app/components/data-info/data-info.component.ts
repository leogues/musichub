import { Component, effect, input } from "@angular/core";

@Component({
  selector: "app-data-info",
  standalone: true,
  imports: [],
  templateUrl: "./data-info.component.html",
})
export class DataInfoComponent {
  isLoading = input<boolean>(true);
  total = input.required<number>();
  selectedCount = input.required<number>();
  filteredCount = input.required<number>();
}
